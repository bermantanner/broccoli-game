package main

// game defines the interface that any playable game or the lobby must implement.
type Game interface {
	// Start is called when the game transitions into this state.
	Start(hub *Hub)

	// HandleMessage processes data sent by a specific client during this game
	HandleMessage(client *Client, msg []byte)

	// AddPlayer is called if someone joins or reconnects mid game
	AddPlayer(client *Client)

	// RemovePlayer is called if someone disconnects
	RemovePlayer(client *Client)

	// GetState returns the current snapshot of the game.
	// used to seamlessly drop reconnecting players back into the game
	GetState() []byte
}
