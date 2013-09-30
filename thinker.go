package checkerlution

type Thinker interface {
	Start(game Game)
	Think(gameState GameState) ValidMove
	Stop()
}
