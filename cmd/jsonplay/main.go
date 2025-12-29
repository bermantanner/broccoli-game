package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Type string `json:"type"`
	Room string `json:"room"`
	Name string `json:"name"`
}

func main() {

	msg := Message{
		Type: "join_room",
		Room: "AXD2",
		Name: "Tanner",
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonMsg))

	var decoded Message
	err = json.Unmarshal(jsonMsg, &decoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(msg)

}
