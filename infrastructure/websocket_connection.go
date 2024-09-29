package infrastructure

import "github.com/gorilla/websocket"

type WebsocketConnection struct {
	conn *websocket.Conn
}

func NewWebsocketConnection(ws *websocket.Conn) *WebsocketConnection {
	return &WebsocketConnection{
		conn: ws,
	}
}

func (ws *WebsocketConnection) ReadMessage() ([]byte, error) {
	_, msg, err := ws.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (ws *WebsocketConnection) WriteMessage(msg []byte) error {
	return ws.conn.WriteMessage(websocket.TextMessage, msg)
}

func (ws *WebsocketConnection) Close() error {
	return ws.conn.Close()
}

func (ws *WebsocketConnection) GetRemoteAddress() string {
	return ws.conn.RemoteAddr().String()
}
