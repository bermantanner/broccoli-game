package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients      map[*Client]bool
	register     chan *Client
	unregister   chan *Client
	broadcast    chan *ClientMessage
	CurrentState Game
}

type ClientMessage struct {
	client *Client
	data   []byte
}

func newHub() *Hub {
	hub := Hub{}
	hub.clients = make(map[*Client]bool)
	hub.register = make(chan *Client)
	hub.unregister = make(chan *Client)
	hub.broadcast = make(chan *ClientMessage)
	// CurrentState starts as nil until we explicitly set it to lobby or game
	return &hub
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			log.Println("registered new client --> total: ", len(h.clients))

			// let the CurrentState handle the new player
			if h.CurrentState != nil {
				h.CurrentState.AddPlayer(c)
			}

		case c := <-h.unregister:
			if _, exists := h.clients[c]; exists {
				delete(h.clients, c)
				close(c.send)
				log.Println("unregistered client --> total: ", len(h.clients))

				// let the CurrentState handle the disconnected player
				if h.CurrentState != nil {
					h.CurrentState.RemovePlayer(c)
				}
			}

		case msgObj := <-h.broadcast:
			if h.CurrentState != nil {
				h.CurrentState.HandleMessage(msgObj.client, msgObj.data)
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

// sends a targeted message to a single specific player in the Hub
func (h *Hub) SendToPlayer(playerName string, msg []byte) {
	for client := range h.clients {
		if client.name == playerName {
			select {
			case client.send <- msg:
			default:
				// if their channel is blocked, drop them
				delete(h.clients, client)
				close(client.send)
			}
			break // stop looping because player was found and msg sent
		}
	}
}
