package main

import (
	"example/websocket/domain"
	"log"
	"net/http"
)

func main() {
	q := domain.NewQueue()
	http.HandleFunc("/ws", q.HandleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
