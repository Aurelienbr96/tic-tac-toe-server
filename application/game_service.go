package application

import (
	"encoding/json"
	"example/websocket/domain"
	"example/websocket/interfaces"
	"fmt"
	"log"
	"sync"
)

type GameService struct {
	Game        interfaces.Game
	Players     *[2]interfaces.Connection
	broadcaster interfaces.Broadcaster
	mu          sync.Mutex
}

func NewGameService(players *[2]interfaces.Connection, broadcaster interfaces.Broadcaster) *GameService {
	return &GameService{
		Players:     players,
		broadcaster: broadcaster,
	}
}

func (gs *GameService) StartGame() {
	gs.Game = domain.NewGame(*gs.Players)

	for i, player := range gs.Players {
		var p string
		if i == 0 {
			p = "x"
		} else {
			p = "o"
		}
		playerUpdate := domain.SetPlayerMessage{
			Player: p,
		}
		gs.broadcaster.SendToPlayer(player, "set-player", playerUpdate)
		go gs.HandlePlayer(player)
	}
}

func (gs *GameService) HandlePlayer(player interfaces.Connection) {
	defer func() {
		gs.mu.Lock()
		player.Close()
		gs.mu.Unlock()
	}()

	fmt.Printf("%s", player.GetRemoteAddress())

	for {
		msg, err := player.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		var receivedMessage domain.ReceivedMessage
		// Print received message
		fmt.Printf("Received: %s\n", msg)

		if err := json.Unmarshal(msg, &receivedMessage); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		if receivedMessage.M == "reset-board" {
			gs.Game.ResetGame()
			gs.broadcaster.Broadcast(gs.Players[:], "reset-board", nil)
		} else {
			turn, _ := gs.Game.SetNextMove(receivedMessage.X, receivedMessage.Y, receivedMessage.M)
			winner := gs.Game.GetWinner()
			if winner == "none" {
				board := gs.Game.GetBoard()
				update := domain.SetBoardUpdateMessage{
					Board: board,
					Turn:  turn,
				}
				gs.broadcaster.Broadcast(gs.Players[:], "board-update", update)
			} else if winner != "none" {
				winnerUpdate := domain.SetWinnerMessage{
					Winner: winner,
				}
				gs.broadcaster.Broadcast(gs.Players[:], "set-winner", winnerUpdate)
			}
		}
		// Print received message
		fmt.Printf("Received: %s\n", msg)
	}
}
