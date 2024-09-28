package infrastructure

import (
	"encoding/json"
	"example/websocket/interfaces"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

type WebsocketBroadcaster struct {
	mu      sync.Mutex
	players map[int]interfaces.Connection
}

func NewWebsocketBroadcaster(playerID int, conn *websocket.Conn) *WebsocketBroadcaster {
	return &WebsocketBroadcaster{
		players: make(map[int]interfaces.Connection),
	}
}

func (broadcast *WebsocketBroadcaster) Broadcast(messageType string, data interface{}) error {
	broadcast.mu.Lock()
	defer broadcast.mu.Unlock()

	var message Message

	msg, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	for i, conn := range broadcast.players {
		if err := conn.WriteMessage(msg); err != nil {
			delete(broadcast.players, i)
			conn.Close()
		}
	}

	return nil
}

func (broadcast *WebsocketBroadcaster) SendToPlayer(playerID int, messageType string, data interface{}) {
	broadcast.mu.Lock()
	defer broadcast.mu.Unlock()

	var message Message

	msg, err := json.Marshal(&message)
	if err != nil {
		log.Printf("error marshaling json %v", err)
	}

	broadcast.players[playerID].WriteMessage(msg)
}

func (broadcast *WebsocketBroadcaster) RemovePlayer(playerID int) error {
	broadcast.mu.Lock()
	defer broadcast.mu.Unlock()

	conn, exists := broadcast.players[playerID]
	if !exists {
		return fmt.Errorf("player %d not connected", playerID)
	}
	conn.Close()
	delete(broadcast.players, playerID)
	return nil
}
