package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	room     string
	name     string
	isHost   bool
	isPlayer bool
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

		role := req.URL.Query().Get("role")
		if role == "" {
			http.Error(w, "missing role in URL param", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("error: ", err)
			return
		}

		client := &Client{
			conn:     conn,
			send:     make(chan []byte, 256),
			room:     room,
			name:     name,
			isHost:   role == "host",
			isPlayer: role != "host",
		}
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

		clientMsg := &ClientMessage{
			client: client,
			data:   msg,
		}
		hub.broadcast <- clientMsg
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
