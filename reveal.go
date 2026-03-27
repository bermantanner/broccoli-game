package main

import (
	"encoding/json"
	"log"
	"time"
)

type Reveal struct {
	hub           *Hub
	Team1         []string
	Team2         []string
	FinalDrawings map[string]string
	TotalRounds   int
	CurrentRound  int
	Duration      int
	Assignments   map[string]string
}

func (r *Reveal) Start(hub *Hub) {
	r.hub = hub

	// broadcast the base64 images to the Host screen
	r.hub.broadcastToAll(r.GetState())

	// for now, start a 20-second timer in the background (update later on)
	go func() {
		time.Sleep(20 * time.Second)

		if r.CurrentRound < r.TotalRounds {
			log.Printf("Starting Round %d!", r.CurrentRound+1)
			nextRound := &Drawing{
				Team1:        r.Team1,
				Team2:        r.Team2,
				TotalRounds:  r.TotalRounds,
				CurrentRound: r.CurrentRound + 1,
				Duration:     r.Duration,
			}
			r.hub.CurrentState = nextRound
			nextRound.Start(r.hub)
		} else {
			// TODO: currently, this doesnt seem to return to lobby, fix later
			// on Go it seems succesful, but frontend does not recieve information
			log.Println("Game Over! Returning to Lobby.")
			lobbyState := &Lobby{}
			r.hub.CurrentState = lobbyState
			lobbyState.Start(r.hub)
		}
	}()
}

func (r *Reveal) HandleMessage(client *Client, msg []byte) {
	// For now, ignore all button clicks during the reveal phase
}

func (r *Reveal) AddPlayer(client *Client) {
	client.send <- r.GetState()
}

func (r *Reveal) RemovePlayer(client *Client) {}

func (r *Reveal) GetState() []byte {
	event := map[string]interface{}{
		"type":         "reveal_started",
		"team1":        r.Team1,
		"team2":        r.Team2,
		"assignments":  r.Assignments,
		"drawings":     r.FinalDrawings,
		"currentRound": r.CurrentRound,
		"totalRounds":  r.TotalRounds,
	}
	b, _ := json.Marshal(event)
	return b
}
