package interfaces

type Broadcaster interface {
	Broadcast(players []Connection, messageType string, data interface{}) error
	SendToPlayer(player Connection, messageType string, data interface{}) error
}
