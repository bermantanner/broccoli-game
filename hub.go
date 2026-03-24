package main

import (
	"encoding/json"
	"fmt"
	"log"
)

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
			h.broadcastToAll(h.lobbySnapshot(c.room))
		case c := <-h.unregister:
			if _, exists := h.clients[c]; exists {
				delete(h.clients, c)
				close(c.send) // this tells writer loop to exit
				log.Println("unregistered client --> total: ", len(h.clients))

				leaveMsg := []byte(fmt.Sprintf(`{"type":"leave","name":%q}`, c.name))
				h.broadcastToAll(leaveMsg)
				h.broadcastToAll(h.lobbySnapshot(c.room))
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

func (h *Hub) lobbySnapshot(room string) []byte {
	players := []string{}
	hostName := ""

	for c := range h.clients {
		if c.isPlayer == true {
			players = append(players, c.name)
		} else {
			hostName = c.name
		}
	}

	event := map[string]any{
		"type":    "lobby",
		"room":    room,
		"players": players,
		"host":    hostName,
	}

	b, _ := json.Marshal(event)
	return b

}
