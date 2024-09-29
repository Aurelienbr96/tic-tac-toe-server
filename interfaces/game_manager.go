package interfaces

type GameManager interface {
	UnqueueTwoPlayers() [2]Connection
	QueuePlayer(player Connection)
	HandleNewPlayer(player Connection)
}
