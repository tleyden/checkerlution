package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	"testing"
)

func init() {
	logg.LogKeys["TEST"] = true
}

func FakePiece() cbot.Piece {
	jsonString := `{"location":7,"validMoves":[{"locations":[1],"captures":[{"team":1,"piece":11}],"king":true}]}`

	piecePtr := &cbot.Piece{}
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, piecePtr)
	if err != nil {
		logg.LogError(err)
	}
	piece := *piecePtr
	return piece

}

func TestNewValidMoveCortexInput(t *testing.T) {

	piece := FakePiece()
	validMove := piece.ValidMoves[0]

	vmCortexInput := NewValidMoveCortexInput(validMove, piece)
	endLocationNormalized := vmCortexInput.endLocation
	expected := -1.0
	logg.LogTo("TEST", "endLocationNormalized: %v", endLocationNormalized)

	assert.True(t, ng.EqualsWithMaxDelta(endLocationNormalized, expected, 0.1))
	assert.True(t, ng.EqualsWithMaxDelta(vmCortexInput.willBecomeKing, 1.0, 0.1))
	assert.True(t, ng.EqualsWithMaxDelta(vmCortexInput.captureValue, 0.0, 0.1))

	logg.LogTo("TEST", "vmCortexInput: %v", vmCortexInput)

}
