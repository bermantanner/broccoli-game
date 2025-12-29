package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func health(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s", req.Method, req.URL.Path)
	fmt.Fprintf(w, "ok\n")
}

func ws(w http.ResponseWriter, req *http.Request) {
	log.Printf("ws connected")
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		msgType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			// this means client disconnected
			log.Printf("ws disconnected")
			break
		}

		log.Printf("recv: %s", string(msgBytes))
		log.Printf("msgType: %d", msgType)

		err = conn.WriteMessage(msgType, msgBytes)
		if err != nil {
			log.Println("error:", err)
			break
		}
	}

}

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/ws", ws)

	http.ListenAndServe(":8090", nil)
}
