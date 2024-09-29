package domain

type Message struct {
	Board      [3][3]string `json:"board"`
	NextPlayer int          `json:"nextPlayer"`
	Type       string       `json:"type"`
	Winner     string       `json:"winner"`
}

type ReceivedMessage struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	M string `json:"m"`
}

type SetBoardUpdateMessage struct {
	Board [3][3]string `json:"board"` // Or any appropriate type
	Turn  string       `json:"turn"`
}

type SetWinnerMessage struct {
	Winner string `json:"winner"`
}

type SetPlayerMessage struct {
	Player string `json:"player"`
}
