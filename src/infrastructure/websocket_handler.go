package infrastructure

import (
	"example/websocket/src/application"
	"example/websocket/src/interfaces"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketHandler struct {
	GameManager interfaces.GameManager
	Broadcaster interfaces.Broadcaster
}

func NewWebsocketHandler(GameManager *application.GameManager, broadcaster interfaces.Broadcaster) *WebsocketHandler {
	return &WebsocketHandler{
		GameManager: GameManager,
		Broadcaster: broadcaster,
	}
}

func (wsHandler WebsocketHandler) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	ws, _ := upgrader.Upgrade(w, r, nil)
	con := NewWebsocketConnection(ws)
	wsHandler.GameManager.HandleNewPlayer(con)
}
