package main

import (
	"log"
	"math/rand"

	"github.com/gorilla/websocket"
)

func GenerateRoomCode() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ConnectToProxy dials the cloud proxy and bridges it to the local Hub
func ConnectToProxy(hub *Hub, proxyURL string) {
	roomCode := GenerateRoomCode()

	// The URL we are dialing (we will use localhost:8080 to test locally first)
	dialURL := proxyURL + "?room=" + roomCode + "&role=engine"

	log.Printf("Connecting to Cloud Proxy to establish Room: %s...", roomCode)

	conn, _, err := websocket.DefaultDialer.Dial(dialURL, nil)
	if err != nil {
		log.Fatal("Fatal error connecting to proxy:", err)
	}

	log.Printf("SUCCESS! Tunnel established. Tell players to join Room: %s", roomCode)

	// Listen for messages coming down the tunnel from the proxy
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Disconnected from proxy:", err)
				break
			}

			// Feed the tunneled message directly into the local state machine!
			// (We will need to slightly adjust how the Hub handles this next)
			log.Printf("Received tunneled message: %s", string(msg))
		}
	}()
}
