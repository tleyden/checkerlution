package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	core "github.com/tleyden/checkers-core"
	"testing"
)

func TestLoadFromBoard(t *testing.T) {

	currentBoardStr := "" +
		"|- o - o - o - o|" +
		"|o - o - o - o -|" +
		"|- - - o - O - o|" +
		"|- - - - x - - -|" +
		"|- - - - - - - -|" +
		"|x - x - o - x -|" +
		"|- x - x - x - x|" +
		"|x - x - x - x -|"

	board := core.NewBoard(currentBoardStr)
	gameStateVector := NewGameStateVector()
	gameStateVector.loadFromBoard(board, core.BLACK_PLAYER)
	assert.Equals(t, gameStateVector[0], OUR_PIECE)
	assert.Equals(t, gameStateVector[31], OPPONENT_PIECE)

}
