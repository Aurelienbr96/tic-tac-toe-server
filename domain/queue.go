package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Queue struct {
	Players    []*Player
	QueueMutex sync.Mutex
}

func NewQueue() *Queue {
	q := &Queue{
		Players: make([]*Player, 0),
	}
	return q
}

func (q *Queue) DequeueTwoPlayers() (*Player, *Player, error) {
	q.QueueMutex.Lock()
	defer q.QueueMutex.Unlock()

	if len(q.Players) < 2 {
		return nil, nil, fmt.Errorf("not enough players in the queue")
	}

	player1 := q.Players[0]
	player2 := q.Players[1]

	// Remove the first two players from the queue
	q.Players = q.Players[2:]

	return player1, player2, nil
}

func (q *Queue) QueuePlayer(p *Player) {
	q.QueueMutex.Lock()
	q.Players = append(q.Players, p)
	q.QueueMutex.Unlock()
}

func (q *Queue) getPlayers() int {
	return len(q.Players)
}

func (q *Queue) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	p := NewPlayer(ws)
	var g *Game

	q.QueuePlayer(p)

	log.Printf("New connection from %s", ws.RemoteAddr())

	if q.getPlayers() == 1 {
		newMsg := InformationMessage{
			Type:    "waiting_room",
			Message: "We are looking for a player for you, please wait...",
		}
		response, _ := json.Marshal(newMsg)
		ws.WriteMessage(websocket.TextMessage, response)
	}

	for q.getPlayers() >= 2 {

		log.Printf("Queue: %v", q)
		player1, player2, _ := q.DequeueTwoPlayers()
		g = NewGame(player1, player2, q)
		msg, err := g.startGame()
		newMsg := InformationMessage{
			Type:    "information",
			Message: msg,
		}
		response, _ := json.Marshal(newMsg)
		ws.WriteMessage(websocket.TextMessage, response)
		if err != nil {
			ws.Close()
		}
	}

}
