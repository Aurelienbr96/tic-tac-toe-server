package main

import (
	"example/websocket/application"
	"example/websocket/infrastructure"
	"log"
	"net/http"
)

// WS handler
// GS BR
// GM handler

func main() {
	broadcaster := infrastructure.NewWebsocketBroadcaster()
	gameManager := application.NewGameManager(broadcaster)
	wsHandler := infrastructure.NewWebsocketHandler(gameManager, broadcaster)

	http.HandleFunc("/ws", wsHandler.HandleNewConnection)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
