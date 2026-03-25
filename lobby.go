package main

import (
	"encoding/json"
	"log"
)

// Lobby represents the "Waiting Room" phase. It implements the Game interface.
type Lobby struct {
	hub *Hub
}

// Start initializes the lobby state
func (l *Lobby) Start(hub *Hub) {
	l.hub = hub
	log.Println("Hub successfully transitioned to Lobby state.")
}

// HandleMessage listens for the Host to click "Start Game"
func (l *Lobby) HandleMessage(client *Client, msg []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("Error parsing JSON in Lobby:", err)
		return
	}

	// If the Host clicks the start button
	if data["type"] == "start_game" && client.isHost {
		log.Println("Host is starting the game!")

		// this is where state will swap later
		l.hub.broadcastToAll(msg)
	}
}

// AddPlayer announces a new arrival
func (l *Lobby) AddPlayer(client *Client) {
	log.Printf("%s joined the lobby.", client.name)
	// broadcast the new player list to everyone
	l.hub.broadcastToAll(l.GetState())
}

// RemovePlayer announces a departure
func (l *Lobby) RemovePlayer(client *Client) {
	log.Printf("%s left the lobby.", client.name)
	// broadcast the new player list to everyone
	l.hub.broadcastToAll(l.GetState())
}

// GetState generates the JSON snapshot of who is currently in the room
func (l *Lobby) GetState() []byte {
	players := []string{}
	hostName := ""

	for c := range l.hub.clients {
		if c.isPlayer {
			players = append(players, c.name)
		} else {
			hostName = c.name
		}
	}

	event := map[string]interface{}{
		"type":    "lobby",
		"players": players,
		"host":    hostName,
	}

	b, _ := json.Marshal(event)
	return b
}
