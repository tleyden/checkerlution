package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"testing"
)

func TestCreateNeurgoCortex(t *testing.T) {
	game := &Game{}
	game.CreateNeurgoCortex()
	cortex := game.cortex
	assert.True(t, cortex != nil)
	assert.True(t, cortex.Sensors != nil)

	cortex.RenderSVGFile("out.svg")

}

func TestChooseBestMove(t *testing.T) {

	ng.SeedRandom()
	logg.LogKeys["MAIN"] = true

	game := &Game{}
	game.CreateNeurgoCortex()
	cortex := game.cortex
	cortex.Run()

	gameState, possibleMoves := game.FetchNewGameDocument()
	bestMove := game.ChooseBestMove(cortex, gameState, possibleMoves)

	found := false
	for _, possibleMove := range possibleMoves {
		if possibleMove == bestMove {
			found = true
		}
	}
	assert.True(t, found)

	cortex.Shutdown()

}

func DISTestGameLoop(t *testing.T) {
	ng.SeedRandom()
	logg.LogKeys["MAIN"] = true
	game := &Game{}
	game.GameLoop()

}
