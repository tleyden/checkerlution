package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestExtractPossibleMoves(t *testing.T) {
	jsonString := FakeGameJson()

	gameState := NewGameStateFromString(jsonString)

	game := &Game{ourTeamId: 0}

	possibleMoves := game.extractPossibleMoves(gameState)

	possibleMove := possibleMoves[0]

	assert.Equals(t, possibleMove.validMove.StartLocation, 7)
	assert.Equals(t, possibleMove.validMove.PieceId, 6)
	assert.Equals(t, len(possibleMoves), 8)

}
