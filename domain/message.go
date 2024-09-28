package domain

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Message struct {
	Board      [3][3]string `json:"board"`
	NextPlayer int          `json:"nextPlayer"`
	Type       string       `json:"type"`
	Winner     string       `json:"winner"`
}

type SetPlayerMessage struct {
	Player string `json:"player"`
	Type   string `json:"type"`
}

type ReceivedMessage struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	M string `json:"m"`
}

type InformationMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func marshalJson(msg any) []byte {
	msgByte, jsonErr := json.Marshal(msg)
	if jsonErr != nil {
		log.Printf("Error marshaling json %v", msg)
	}
	return msgByte
}

func SendMessage(msg any, ws *websocket.Conn) {
	msgByte := marshalJson(msg)
	err := ws.WriteMessage(websocket.TextMessage, msgByte)
	if err != nil {
		log.Printf("Error sending message to player %v", ws)
	}
}
