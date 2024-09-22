package domain

type Message struct {
	Board      [3][3]string `json:"board"`
	NextPlayer int          `json:"nextPlayer"`
	Type       string       `json:"type"`
	Winner     string       `json:"winner"`
}

type SetPlayerMessage struct {
	Player string `json:"player"`
	Type   string `json:"type"`
}

type ReceivedMessage struct {
	X int    `json:"x"`
	Y int    `json:"y"`
	M string `json:"m"`
}

type InformationMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
