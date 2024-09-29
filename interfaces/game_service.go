package interfaces

type GameService interface {
	StartGame()
	HandlePlayer(player Connection)
}
