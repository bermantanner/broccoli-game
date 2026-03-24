package main

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

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

type CreateRoomResponse struct {
	Room string `json:"room"`
}
