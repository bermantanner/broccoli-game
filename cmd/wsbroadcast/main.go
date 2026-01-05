package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	room   string
	name   string
	isHost bool
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

			joinMsg := []byte(fmt.Sprintf(`{"type":"join","name":%q}`, c.name))
			h.broadcastToAll(joinMsg)
		case c := <-h.unregister:
			if _, exists := h.clients[c]; exists {
				delete(h.clients, c)
				close(c.send) // this tells writer loop to exit
				log.Println("unregistered client --> total: ", len(h.clients))

				leaveMsg := []byte(fmt.Sprintf(`{"type":"leave","name":%q}`, c.name))
				h.broadcastToAll(leaveMsg)
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

func (h *Hub) broadcastToAll(msg []byte) {
	for c := range h.clients {
		select {
		case c.send <- msg:
		default:
			delete(h.clients, c)
			close(c.send)
		}
	}
}

type RoomManager struct {
	mu    sync.Mutex
	rooms map[string]*Hub
}

func newRoomManager() *RoomManager {
	roomManager := RoomManager{}
	roomManager.rooms = make(map[string]*Hub)

	return &roomManager
}

func (rm *RoomManager) GetOrCreateRoom(code string) *Hub {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	hub, exists := rm.rooms[code]

	if exists {
		return hub
	}

	hub = newHub()
	go hub.run()
	rm.rooms[code] = hub
	return hub
}

func (rm *RoomManager) GetRoom(code string) (*Hub, bool) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	hub, exists := rm.rooms[code]
	return hub, exists
}

func createRoomHandler(rm *RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		code := generateRoomCode()
		rm.GetOrCreateRoom(code)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateRoomResponse{Room: code})
	}
}

func generateRoomCode() string {
	//later on we should make it regenerate if a code exists.
	const length = 4
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	code := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			panic(err)
		}
		code[i] = alphabet[n.Int64()]
	}

	return string(code)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(rm *RoomManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		room := req.URL.Query().Get("room")
		if room == "" {
			http.Error(w, "missing room code in URL param", http.StatusBadRequest)
			return
		}

		hub, exists := rm.GetRoom(room)
		if !exists {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		name := req.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "missing name in URL param", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("error: ", err)
			return
		}

		client := &Client{conn: conn, send: make(chan []byte, 256), room: room, name: name}
		hub.register <- client

		go writePump(client)
		readPump(hub, client)

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

type CreateRoomResponse struct {
	Room string `json:"room"`
}

type ServerEvent struct {
	Type    string   `json:"type"` // "join", "leave", "chat", "lobby"
	Room    string   `json:"room"`
	Name    string   `json:"name,omitempty"`
	Host    bool     `json:"host,omitempty"`
	Text    string   `json:"text,omitempty"`
	Players []string `json:"players,omitempty"`
}

func main() {
	roomManager := newRoomManager()

	http.HandleFunc("/create-room", createRoomHandler(roomManager))
	http.HandleFunc("/ws", wsHandler(roomManager))
	log.Fatal(http.ListenAndServe(":8090", nil))
}
