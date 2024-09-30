package main

import (
	"encoding/json"
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
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			myMsg := map[string]string{
				"Hello": "Server Health OK",
			}
			msg, _ := json.Marshal(myMsg)
			w.Write(msg)
		}
	})

	log.Println("Server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
