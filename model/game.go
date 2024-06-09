package model

type GameState struct {
	GameID string `json:"gameId"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Score  int    `json:"score"`
	Fruit  Fruit  `json:"fruit"`
	Snake  Snake  `json:"snake"`
}

type Fruit struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Snake struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	VelX int `json:"velX"` // X velocity of the snake (one of -1, 0, 1)
	VelY int `json:"velY"` // Y velocity of the snake (one of -1, 0, 1)
}
