package domain

import (
	"example/websocket/src/interfaces"
	"fmt"
)

var winningCombinations = [][][2]int{
	// rows
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// cols
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	// diag
	{{0, 0}, {1, 1}, {2, 2}},
	{{2, 0}, {1, 1}, {0, 2}},
}

var initialState = [3][3]string{}

type Game struct {
	Board [3][3]string
	Turn  string
}

func NewGame(players [2]interfaces.Connection) *Game {
	return &Game{
		Board: initialState,
		Turn:  "x",
	}
}

func (g *Game) IsBoardFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func (g *Game) ResetGame() [3][3]string {
	g.Board = initialState
	g.Turn = "x"
	return g.Board
}

func (g *Game) GetWinner() string {
	for _, combination := range winningCombinations {
		cell1 := g.Board[combination[0][0]][combination[0][1]]
		cell2 := g.Board[combination[1][0]][combination[1][1]]
		cell3 := g.Board[combination[2][0]][combination[2][1]]

		if cell1 != "" && cell1 == cell2 && cell2 == cell3 {
			return cell1
		}
	}

	if g.IsBoardFull() {
		return "draw"
	}

	return "none"
}

func (g *Game) SetNextMove(x int, y int, m string) (string, error) {
	if x > 2 || y > 2 || y < 0 || x < 0 {
		return "", fmt.Errorf("invalid board position")
	}
	if g.Board[x][y] != "" {
		return "", fmt.Errorf("invalid board position")
	}
	g.Board[x][y] = m

	// set next player

	if g.Turn == "x" {
		g.Turn = "o"
	} else {
		g.Turn = "x"
	}

	return g.Turn, nil
}

func (g *Game) GetBoard() [3][3]string {
	return g.Board
}
