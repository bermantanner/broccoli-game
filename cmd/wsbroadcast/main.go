package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newHub() *Hub {
	hub := Hub{}
	hub.clients = make(map[*Client]bool)
	hub.register = make(chan *Client)
	hub.unregister = make(chan *Client)
	hub.broadcast = make(chan []byte)

	return &hub
}

func (h *Hub) run() {
	for {
		select { // this is basically "pause here until one of these channels has data"
		case c := <-h.register:
			h.clients[c] = true
			log.Println("registered new client --> total: ", len(h.clients))
		case c := <-h.unregister:
			if _, exists := h.clients[c]; exists {
				delete(h.clients, c)
				close(c.send) // this tells writer loop to exit
				log.Println("unregistered client --> total: ", len(h.clients))
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// client's send channel is full or blocked, so we drop client
					delete(h.clients, c)
					close(c.send)
				}

			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// this is a bit confusing, but the traditional handlerfunc is unable to know
// about our hub unless we make it public (which is bad practice)
// so, we must make this handler function that knows about our hub
func wsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("error: ", err)
			return
		}

		client := Client{conn: conn, send: make(chan []byte, 256)}
		hub.register <- &client

		go writePump(&client)
		readPump(hub, &client)

	}
}

// socket sending message to server
func readPump(hub *Hub, client *Client) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
	}()

	for {
		msgType, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("error reading msg from client: ", err)
			break
		}
		// handling just txt for now
		if msgType != websocket.TextMessage {
			continue
		}

		log.Printf("recieved: %s", string(msg))

		hub.broadcast <- msg
	}
}

// server sending message to socket
func writePump(client *Client) {
	defer client.conn.Close()

	for msg := range client.send { // this is blocking! channels block when empty.
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/ws", wsHandler(hub))
	http.ListenAndServe(":8090", nil)
}
