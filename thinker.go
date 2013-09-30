package checkerlution

type Thinker interface {
	Start(ourTeamId int)
	Think(gameState GameState) ValidMove
	Stop()
}
