package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var winningCombinations = [][][2]int{
	// rows
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// cols
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	// diag
	{{0, 0}, {1, 1}, {2, 2}},
	{{2, 0}, {1, 1}, {0, 2}},
}

var initialState = [3][3]string{}

type Game struct {
	Board     [3][3]string
	Players   []*Player
	Turn      int
	GameMutex sync.Mutex
	Queue     *Queue
}

func NewGame(player1, player2 *Player, q *Queue) *Game {
	return &Game{
		Board:   initialState,
		Players: []*Player{player1, player2},
		Turn:    0,
		Queue:   q,
	}
}

func (g *Game) isBoardFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func (g *Game) broadcastBoard() {
	winner := GetWinner(g)
	log.Printf("%s", winner)

	message := Message{
		Type:       "board_update",
		Board:      g.Board,
		Winner:     winner,
		NextPlayer: g.Turn,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
	}

	for i, player := range g.Players {
		if player.websockCon != nil {
			err := player.websockCon.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				log.Printf("Error sending message to player %d: %v", i, err)
				g.GameMutex.Lock()
				g.Players[i].websockCon = nil
				g.GameMutex.Unlock()
				player.websockCon.Close()
			} else {
				log.Printf("Message sent to player %d", i)
			}
		}
	}
}

func (g *Game) resetBoard() {
	g.Board = initialState
}

func (g *Game) cleanupGame() {
	g.resetBoard()
	// Optional: clear player list or reset game state completely
	g.Players = nil
	g = nil
}

func (g *Game) handleClient(ws *websocket.Conn, playerIndex int) {
	defer func() {

		g.GameMutex.Lock()
		log.Printf("game player length: %d, game: %v", len(g.Players), g)

		g.Players[playerIndex].websockCon = nil
		if playerIndex == 1 && g.Players[0] != nil {
			g.Queue.QueuePlayer(g.Players[0])

			sendMessagePlayer := SetPlayerMessage{
				Type: "set-player",
			}

			msgByte, _ := json.Marshal(sendMessagePlayer)
			g.Players[0].websockCon.WriteMessage(websocket.TextMessage, msgByte)
		} else if g.Players[1] != nil {
			g.Queue.QueuePlayer(g.Players[1])

			sendMessagePlayer := SetPlayerMessage{
				Type: "set-player",
			}

			msgByte, _ := json.Marshal(sendMessagePlayer)
			g.Players[1].websockCon.WriteMessage(websocket.TextMessage, msgByte)
		}
		g.GameMutex.Unlock()
		// close the game
		g.cleanupGame()
		ws.Close()
		log.Printf("Client %v disconnected", ws.RemoteAddr())
	}()

	fmt.Printf("handle client %s", ws.RemoteAddr())

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		var receivedMessage ReceivedMessage
		// Print received message
		fmt.Printf("Received: %s\n", msg)

		err = json.Unmarshal(msg, &receivedMessage)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			continue
		}
		if receivedMessage.M == "reset_board" {
			g.resetBoard()
		} else {

			SetNextMove(receivedMessage.X, receivedMessage.Y, receivedMessage.M, g)
		}
		g.broadcastBoard()
	}
}

func (g *Game) startGame() (string, error) {
	if len(g.Players) < 2 {
		return "not enough player in the game", fmt.Errorf("not enough player in the game")
	}

	for i := 0; i < len(g.Players); i++ {
		log.Printf("player: %d", g.Players[i].websockCon.RemoteAddr())

		if i == 0 {
			sendMessagePlayer := SetPlayerMessage{
				Player: "x",
				Type:   "set-player",
			}
			msgByte, _ := json.Marshal(sendMessagePlayer)
			g.Players[i].websockCon.WriteMessage(websocket.TextMessage, msgByte)
		} else {
			sendMessagePlayer := SetPlayerMessage{
				Player: "o",
				Type:   "set-player",
			}
			msgByte, _ := json.Marshal(sendMessagePlayer)
			g.Players[i].websockCon.WriteMessage(websocket.TextMessage, msgByte)
		}
		go g.handleClient(g.Players[i].websockCon, i)
	}

	connectedPlayers := GetConnectedPlayers(g)
	log.Printf("Connected players: %d", connectedPlayers)

	// Start the client handler

	return "Game started correctly", nil
}
