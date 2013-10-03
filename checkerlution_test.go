package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	"testing"
)

func TestExtractPossibleMoves(t *testing.T) {
	jsonString := FakeGameJson()

	gameState := cbot.NewGameStateFromString(jsonString)

	checkerlution := &Checkerlution{}
	checkerlution.ourTeamId = cbot.RED_TEAM

	possibleMoves := checkerlution.extractPossibleMoves(gameState)

	possibleMove := possibleMoves[0]

	assert.Equals(t, possibleMove.validMove.StartLocation, 7)
	assert.Equals(t, possibleMove.validMove.PieceId, 6)
	assert.Equals(t, len(possibleMoves), 8)

}

func TestExtractGameStateVector(t *testing.T) {

	jsonString := FakeGameJson()

	gameState := cbot.NewGameStateFromString(jsonString)

	checkerlution := &Checkerlution{}
	checkerlution.ourTeamId = cbot.RED_TEAM

	gameStateVector := checkerlution.extractGameStateVector(gameState)

	logg.LogTo("TEST", "gameStateVector: %v", gameStateVector)
	// at location 1, which has index 0, is our (Team 0) king.
	// so we expect to see a 1.0 there
	assert.True(t, gameStateVector[0] == OUR_KING)

	// the opponent has a normal piece on last square (location 32)
	assert.True(t, gameStateVector[31] == OPPONENT_PIECE)

}

func TestChooseBestMove(t *testing.T) {

	ng.SeedRandom()
	logg.LogKeys["MAIN"] = true

	checkerlution := &Checkerlution{}
	checkerlution.ourTeamId = cbot.RED_TEAM

	checkerlution.CreateNeurgoCortex()
	cortex := checkerlution.cortex
	cortex.Run()

	gameState, possibleMoves := FakeGameDocument()
	bestMove := checkerlution.chooseBestMove(gameState, possibleMoves)
	logg.LogTo("TEST", "bestMove: %v", &bestMove)

	found := false
	for _, possibleMove := range possibleMoves {
		logg.LogTo("TEST", "possibleMove: %v", &possibleMove)
		if possibleMove.Equals(bestMove) {
			found = true
		}
	}
	assert.True(t, found)

	cortex.Shutdown()

}

func FakeGameDocument() (gameState []float64, possibleMoves []ValidMoveCortexInput) {

	gameState = NewGameStateVector()

	possibleMove1 := ValidMoveCortexInput{
		startLocation:   0,
		isCurrentlyKing: -1,
		endLocation:     1.0,
		willBecomeKing:  -0.5,
		captureValue:    1,
	}

	possibleMove2 := ValidMoveCortexInput{
		startLocation:   1,
		isCurrentlyKing: -0.5,
		endLocation:     0.0,
		willBecomeKing:  0.5,
		captureValue:    0,
	}

	possibleMoves = []ValidMoveCortexInput{possibleMove1, possibleMove2}
	return

}

func FakeGameJson() string {
	jsonString := `{"applicationUrl":"http://www.couchbase.com/checkers","applicationName":"Couchbase Checkers","revotingAllowed":false,"highlightPiecesWithMoves":true,"number":1,"startTime":"2013-08-26T16:05:30Z","moveDeadline":"2013-08-26T16:05:45Z","turn":1,"activeTeam":0,"winningTeam":0,"moves":[],"teams":[{"participantCount":117983,"score":11,"pieces":[{"location":1,"king":true},{"location":2},{"location":3},{"location":4},{"location":5},{"location":6},{"location":7,"validMoves":[{"locations":[11],"captures":[{"team":1,"piece":11}],"king":true}]},{"location":8,"validMoves":[{"locations":[11],"captures":[{"team":1,"piece":8},{"team":1,"piece":9},{"team":1,"piece":10}]},{"locations":[11,15]}]},{"location":9,"validMoves":[{"locations":[13]},{"locations":[14]}]},{"location":10,"validMoves":[{"locations":[14]},{"locations":[15]}]},{"location":11,"captured":true},{"location":12,"king":true,"validMoves":[{"locations":[16]}]}]},{"participantCount":109217,"score":12,"pieces":[{"location":21},{"location":22},{"location":23},{"location":24},{"location":25},{"location":26},{"location":27},{"location":28},{"location":29},{"location":30},{"location":31},{"location":32}]}]}`
	return jsonString
}
