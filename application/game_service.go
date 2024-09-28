package application

import (
	"example/websocket/domain"
	"example/websocket/interfaces"
	"fmt"
	"log"
	"sync"
)

type GameService struct {
	Game        *domain.Game
	Players     *[2]interfaces.Connection
	broadcaster *interfaces.Broadcaster
	mu          sync.Mutex
}

func NewGameService(players *[2]interfaces.Connection, broadcaster interfaces.Broadcaster) *GameService {
	return &GameService{
		Players: players,
	}
}

func (gs *GameService) StartGame() {
	gs.Game = domain.NewGame(*gs.Players)
	for _, player := range gs.Players {
		go gs.HandlePlayer(player)
	}
}

func (gs *GameService) HandlePlayer(player interfaces.Connection) {
	defer func() {
		gs.mu.Lock()
		player.Close()
		gs.mu.Unlock()
	}()

	fmt.Printf("handle client %s", player.GetRemoteAddress())

	for {
		msg, err := player.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Print received message
		fmt.Printf("Received: %s\n", msg)

	}
}
