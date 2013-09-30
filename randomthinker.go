package checkerlution

type RandomThinker struct {
	ourTeamId int
}

func (r *RandomThinker) Start(ourTeamId int) {
	r.ourTeamId = ourTeamId
}

func (r RandomThinker) Think(gameState GameState) (bestMove ValidMove) {
	ourTeam := gameState.Teams[r.ourTeamId]
	allValidMoves := ourTeam.AllValidMoves()
	randomValidMoveIndex := randomIntInRange(0, len(allValidMoves))
	bestMove = allValidMoves[randomValidMoveIndex]
	return
}

func (r RandomThinker) Stop() {

}
