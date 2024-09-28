package domain

import (
	"example/websocket/interfaces"
	"fmt"
	"sync"
)

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
	Board [3][3]string
	Turn  string
	mu    sync.Mutex
}

func NewGame(players [2]interfaces.Connection) *Game {
	return &Game{
		Board: initialState,
		Turn:  "X",
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

/* func (g *Game) broadcastBoard() {
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
} */

func (g *Game) resetBoard() {
	g.Board = initialState
}

/* func (g *Game) cleanupGame() {
	g.resetBoard()
	// Optional: clear player list or reset game state completely
	g.Players = nil
	g = nil
} */

/* func (g *Game) handleClient(ws *websocket.Conn, playerIndex int) {
	// stop the game if someone disconnect
	defer func() {
		g.GameMutex.Lock()
		log.Printf("game player length: %d, game: %v", len(g.Players), g)
		for i, player := range g.Players {
			if i != playerIndex {
				sendMessagePlayer := InformationMessage{
					Type:    "close-game",
					Message: "Your opponent disconnected",
				}
				SendMessage(sendMessagePlayer, player.websockCon)
			}
			player.websockCon.Close()
			log.Printf("Client %v disconnected", player.websockCon.RemoteAddr())
		}
		g.cleanupGame()
		g.GameMutex.Unlock()
		// close the game
		// g.cleanupGame()
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
} */

func (g *Game) GetWinner() string {
	for _, combination := range winningCombinations {
		cell1 := g.Board[combination[0][0]][combination[0][1]]
		cell2 := g.Board[combination[1][0]][combination[1][1]]
		cell3 := g.Board[combination[2][0]][combination[2][1]]

		if cell1 != "" && cell1 == cell2 && cell2 == cell3 {
			return cell1
		}
	}

	if g.isBoardFull() {
		return "draw"
	}

	return "none"
}

func (g *Game) SetNextMove(x int, y int, m string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if x > 2 || y > 2 || y < 0 || x < 0 {
		return fmt.Errorf("invalid board position")
	}
	if g.Board[x][y] != "" {
		return fmt.Errorf("invalid board position")
	}
	g.Board[x][y] = m

	// set next player

	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}

	return nil
}

/* func (g *Game) startGame() (string, error) {
	if len(g.Players) < 2 {
		return "not enough player in the game", fmt.Errorf("not enough player in the game")
	}

	for i := 0; i < len(g.Players); i++ {
		log.Printf("player: %d", g.Players[i].websockCon.RemoteAddr())

		var sendMessagePlayer SetPlayerMessage
		// first connected player X
		if i == 0 {
			sendMessagePlayer = SetPlayerMessage{
				Player: "x",
				Type:   "set-player",
			}
			// first connected player O
		} else {
			sendMessagePlayer = SetPlayerMessage{
				Player: "o",
				Type:   "set-player",
			}
		}
		SendMessage(sendMessagePlayer, g.Players[i].websockCon)
		go g.handleClient(g.Players[i].websockCon, i)
	}

	connectedPlayers := GetConnectedPlayers(g)
	log.Printf("Connected players: %d", connectedPlayers)

	// Start the client handler

	return "Game started correctly", nil
} */
