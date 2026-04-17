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

	l.hub.broadcastToAll(l.GetState())
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
		gameName, _ := data["game"].(string)

		// enforce game-specific player limits
		if gameName == "drawn together" {
			playerCount := 0
			for c := range l.hub.clients {
				if c.isPlayer {
					playerCount++
				}
			}

			if playerCount != 4 && playerCount != 6 {
				log.Printf("Start blocked: Drawn Together requires 4 or 6 players. Current: %d", playerCount)

				// send a targeted error message only to Host
				errorMsg := map[string]string{
					"type":    "error",
					"message": "Drawn Together requires exactly 4 or 6 players to start!",
				}
				b, _ := json.Marshal(errorMsg)
				client.send <- b
				return // stop executing! do not transition states.
			}
		}

		log.Println("Host is starting the game! Transitioning to Team Select.")
		// extract the options (JSON numbers parse as float64)
		rounds, _ := data["rounds"].(float64)
		duration, _ := data["duration"].(float64)

		// unplug the Lobby, plug in TeamSelect and pass the settings!
		l.hub.CurrentState = &TeamSelect{
			TotalRounds: int(rounds),
			Duration:    int(duration),
		}
		l.hub.CurrentState.Start(l.hub)
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
