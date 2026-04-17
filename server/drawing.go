package main

import (
	"encoding/json"
	"log"
	"math/rand"
)

// Drawing represents the active game loop where players are on the canvas
type Drawing struct {
	hub           *Hub
	Team1         []string
	Team2         []string
	Assignments   map[string]string
	FinalDrawings map[string]string
	TotalRounds   int
	CurrentRound  int
	Duration      int
}

// Start assigns the canvas sections and broadcasts the start signal
func (d *Drawing) Start(hub *Hub) {
	d.hub = hub
	d.Assignments = make(map[string]string)
	d.FinalDrawings = make(map[string]string)

	// Assign roles for both teams
	d.assignRoles(d.Team1)
	d.assignRoles(d.Team2)

	log.Printf("Drawing phase started! Assignments: %v", d.Assignments)
	d.hub.broadcastToAll(d.GetState())
}

// assignRoles randomly assigns canvas sections based on team size
func (d *Drawing) assignRoles(team []string) {
	shuffled := make([]string, len(team))
	copy(shuffled, team)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	if len(shuffled) == 3 {
		d.Assignments[shuffled[0]] = "head"
		d.Assignments[shuffled[1]] = "body"
		d.Assignments[shuffled[2]] = "legs"
	} else if len(shuffled) == 2 {
		d.Assignments[shuffled[0]] = "top_half"
		d.Assignments[shuffled[1]] = "bottom_half"
	}
}

// HandleMessage acts as the high-speed router and the final collector
func (d *Drawing) HandleMessage(client *Client, msg []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return
	}

	msgType, _ := data["type"].(string)

	// ==========================================
	// routing live strokes
	// ==========================================
	if msgType == "draw_line" && client.isPlayer {
		targetRole, ok := data["targetRole"].(string)
		if !ok {
			return
		}

		// figure out which team the sender is on
		var myTeam []string
		if contains(d.Team1, client.name) {
			myTeam = d.Team1
		} else if contains(d.Team2, client.name) {
			myTeam = d.Team2
		}

		// find the teammate who has the targetRole
		var targetName string
		for _, name := range myTeam {
			if d.Assignments[name] == targetRole {
				targetName = name
				break
			}
		}

		// forward the exact raw byte message straight to them
		if targetName != "" {
			for c := range d.hub.clients {
				if c.name == targetName {
					c.send <- msg
					break
				}
			}
		}
	}
	// time's up.
	if msgType == "time_up" && client.isHost {
		log.Println("Host says time is up!")
		d.hub.broadcastToAll(msg)
	}

	// ==========================================
	// COLLECTING: Catch final images
	// ==========================================
	if msgType == "submit_drawing" && client.isPlayer {
		imageData, ok := data["image"].(string)
		if ok {
			d.FinalDrawings[client.name] = imageData
			log.Printf("Received final drawing from %s", client.name)

			// Check if we have everyone's drawings
			totalPlayers := len(d.Team1) + len(d.Team2)
			if len(d.FinalDrawings) == totalPlayers {
				log.Println("All drawings received! Transitioning to Reveal state...")

				revealState := &Reveal{
					Team1:         d.Team1,
					Team2:         d.Team2,
					Assignments:   d.Assignments,
					FinalDrawings: d.FinalDrawings,
					TotalRounds:   d.TotalRounds,
					CurrentRound:  d.CurrentRound,
					Duration:      d.Duration,
				}
				d.hub.CurrentState = revealState
				revealState.Start(d.hub)
			}
		}
	}
}

func (d *Drawing) AddPlayer(client *Client) {
	client.send <- d.GetState()
}

func (d *Drawing) RemovePlayer(client *Client) {
	log.Printf("%s disconnected during the drawing phase.", client.name)
}

func (d *Drawing) GetState() []byte {
	event := map[string]interface{}{
		"type":         "drawing_started",
		"assignments":  d.Assignments,
		"team1":        d.Team1,
		"team2":        d.Team2,
		"duration":     d.Duration,
		"currentRound": d.CurrentRound,
		"totalRounds":  d.TotalRounds,
	}
	b, _ := json.Marshal(event)
	return b
}

// --- HELPER ---
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
