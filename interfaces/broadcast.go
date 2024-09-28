package interfaces

type Broadcaster interface {
	Broadcast(messageType string, data interface{}) error
	SendToPlayer(playerID int, messageType string, data interface{}) error
	RemovePlayer(playerID int) error
}
