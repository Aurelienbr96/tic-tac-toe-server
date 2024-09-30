package interfaces

type Game interface {
	IsBoardFull() bool
	ResetGame() [3][3]string
	GetWinner() string
	SetNextMove(x int, y int, m string) (string, error)
	GetBoard() [3][3]string
}
