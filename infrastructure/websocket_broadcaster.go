package infrastructure

import (
	"encoding/json"
	"example/websocket/interfaces"
	"log"
)

type Message struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

type WebsocketBroadcaster struct {
}

func NewWebsocketBroadcaster() *WebsocketBroadcaster {
	return &WebsocketBroadcaster{}
}

func (broadcast *WebsocketBroadcaster) Broadcast(players []interfaces.Connection, messageType string, data interface{}) error {
	message := Message{
		Type: messageType,
		Data: data,
	}

	msg, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	for _, conn := range players {
		if err := conn.WriteMessage(msg); err != nil {
			conn.Close()
		}
	}

	return nil
}

func (broadcast *WebsocketBroadcaster) SendToPlayer(player interfaces.Connection, messageType string, data interface{}) error {
	message := Message{
		Data: data,
		Type: messageType,
	}

	msg, err := json.Marshal(&message)
	if err != nil {
		log.Printf("error marshaling json %v", err)
		return err
	}

	player.WriteMessage(msg)
	return nil
}
