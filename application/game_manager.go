package application

import (
	"example/websocket/domain"
	"example/websocket/interfaces"
	"sync"
)

type GameManager struct {
	games       []*domain.Game
	players     []interfaces.Connection
	broadcaster interfaces.Broadcaster
	mu          sync.Mutex
}

func NewGameManager(broadcaster interfaces.Broadcaster) *GameManager {
	return &GameManager{
		broadcaster: broadcaster,
	}
}

func (g *GameManager) UnqueueTwoPlayers() [2]interfaces.Connection {
	g.mu.Lock()
	defer g.mu.Unlock()

	var players [2]interfaces.Connection
	if len(g.players) >= 2 {
		players[0] = g.players[len(g.players)]
		players[1] = g.players[len(g.players)-1]
		g.players = g.players[2:]
	}
	return players
}

func (g *GameManager) QueuePlayer(player interfaces.Connection) {
	g.players = append(g.players, player)
}

func (g *GameManager) HandleNewPlayer(player interfaces.Connection) {
	g.QueuePlayer(player)

	for len(g.players) >= 2 {
		g.mu.Lock()
		players := g.UnqueueTwoPlayers()

		gameService := NewGameService(&players, g.broadcaster)
		gameService.StartGame()

		g.mu.Unlock()
	}
}
