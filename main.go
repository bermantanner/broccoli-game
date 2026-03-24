package main

import (
	"log"
	"net/http"
)

func main() {
	roomManager := newRoomManager()

	http.HandleFunc("/create-room", createRoomHandler(roomManager))
	http.HandleFunc("/ws", wsHandler(roomManager))

	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	log.Println("Server starting on :8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
