package domain

import "github.com/gorilla/websocket"

type Player struct {
	websockCon *websocket.Conn
}

func NewPlayer(player *websocket.Conn) *Player {
	p := &Player{
		websockCon: player,
	}
	return p
}
