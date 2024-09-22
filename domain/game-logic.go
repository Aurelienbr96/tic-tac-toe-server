package domain

import "fmt"

func GetConnectedPlayers(g *Game) int {
	count := 0
	g.GameMutex.Lock()
	defer g.GameMutex.Unlock()

	for _, player := range g.Players {
		if player.websockCon != nil {
			count++
		}
	}
	return count
}

func SetNextMove(x int, y int, m string, g *Game) error {
	g.GameMutex.Lock()
	defer g.GameMutex.Unlock()
	if x > 2 || y > 2 || y < 0 || x < 0 {
		return fmt.Errorf("invalid board position")
	}
	if g.Board[x][y] != "" {
		return fmt.Errorf("invalid board position")
	}
	g.Board[x][y] = m

	// set next player

	if g.Turn == 0 {
		g.Turn = 1
	} else {
		g.Turn = 0
	}

	return nil
}

func GetWinner(g *Game) string {
	for _, combination := range winningCombinations {
		cell1 := g.Board[combination[0][0]][combination[0][1]]
		cell2 := g.Board[combination[1][0]][combination[1][1]]
		cell3 := g.Board[combination[2][0]][combination[2][1]]

		if cell1 != "" && cell1 == cell2 && cell2 == cell3 {
			return cell1
		}
	}

	if g.isBoardFull() {
		return "draw"
	}

	return "none"
}
