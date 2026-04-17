package main

import (
	"encoding/json"
	"log"
)

// TeamSelect represents the phase where players pick Left or Right (Team 1 or Team 2)
// currently this is only used for 'Drawn Together'
type TeamSelect struct {
	hub         *Hub
	Team1       []string
	Team2       []string
	Unassigned  []string
	MaxPerTeam  int
	TotalRounds int
	Duration    int
}

// Start populates the initial Unassigned list and sets max team size
func (ts *TeamSelect) Start(hub *Hub) {
	ts.hub = hub
	ts.Team1 = []string{}
	ts.Team2 = []string{}
	ts.Unassigned = []string{}

	// gather all current players
	playerCount := 0
	for c := range hub.clients {
		if c.isPlayer {
			ts.Unassigned = append(ts.Unassigned, c.name)
			playerCount++
		}
	}

	// since the lobby enforces 4 or 6 players, MaxPerTeam is half
	ts.MaxPerTeam = playerCount / 2

	log.Printf("Team Select Started. Total Players: %d, Max per team: %d", playerCount, ts.MaxPerTeam)
	ts.hub.broadcastToAll(ts.GetState())
}

// HandleMessage listens for players tapping a team, or the host confirming
func (ts *TeamSelect) HandleMessage(client *Client, msg []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("Error parsing JSON in TeamSelect:", err)
		return
	}

	msgType := data["type"]

	// PLAYER ACTION: tapping a team button
	if msgType == "join_team" && client.isPlayer {
		teamWanted := int(data["team"].(float64)) // JSON numbers parse as float64

		// check capacity
		if teamWanted == 1 && len(ts.Team1) >= ts.MaxPerTeam {
			return // Team 1 is full, ignore the tap
		}
		if teamWanted == 2 && len(ts.Team2) >= ts.MaxPerTeam {
			return // Team 2 is full, ignore the tap
		}

		// remove them from wherever they currently are
		ts.removeFromAll(client.name)

		// place them in the new team
		if teamWanted == 1 {
			ts.Team1 = append(ts.Team1, client.name)
		} else if teamWanted == 2 {
			ts.Team2 = append(ts.Team2, client.name)
		}

		// broadcast the updated lists so the UI visually moves their name
		ts.hub.broadcastToAll(ts.GetState())
	}

	// HOST ACTION: clicking "Lock Teams / Next"
	if msgType == "confirm_teams" && client.isHost {
		log.Println("Teams locked! Transitioning to drawing phase...")

		drawingState := &Drawing{
			Team1:        ts.Team1,
			Team2:        ts.Team2,
			TotalRounds:  ts.TotalRounds,
			CurrentRound: 1, // Start at round 1!
			Duration:     ts.Duration,
		}

		ts.hub.CurrentState = drawingState
		ts.hub.CurrentState.Start(ts.hub)
	}
}

// AddPlayer / RemovePlayer handle disconnects gracefully
func (ts *TeamSelect) AddPlayer(client *Client) {
	client.send <- ts.GetState()
}

func (ts *TeamSelect) RemovePlayer(client *Client) {
	log.Printf("%s disconnected during team select.", client.name)
}

// GetState generates the JSON snapshot of the teams
func (ts *TeamSelect) GetState() []byte {
	event := map[string]interface{}{
		"type":       "team_update",
		"team1":      ts.Team1,
		"team2":      ts.Team2,
		"unassigned": ts.Unassigned,
	}
	b, _ := json.Marshal(event)
	return b
}

// helper functions ---

// removeFromAll scrubs the player from all lists so they can switch sides cleanly
func (ts *TeamSelect) removeFromAll(name string) {
	ts.Team1 = removeString(ts.Team1, name)
	ts.Team2 = removeString(ts.Team2, name)
	ts.Unassigned = removeString(ts.Unassigned, name)
}

// removeString is a simple utility to delete a specific string from a slice
func removeString(slice []string, val string) []string {
	var result []string
	for _, v := range slice {
		if v != val {
			result = append(result, v)
		}
	}
	return result
}
