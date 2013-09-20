package checkerlution

import (
	"code.google.com/p/dsallings-couch-go"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"testing"
)

func TestChangesFeed(t *testing.T) {

	logg.LogNoColor()
	logg.LogKeys["TEST"] = true
	logg.LogKeys["DEBUG"] = true
	db, error := couch.Connect(SERVER_URL)
	if error != nil {
		logg.LogPanic("Error connecting to %v: %v", SERVER_URL, error)
	}
	logg.LogTo("TEST", "db: %v", db)

	options := make(map[string]interface{})
	options["since"] = 0
	logg.LogTo("TEST", "calling db.Changes with: %v", handleChange)
	db.Changes(handleChange, options)
	logg.LogTo("TEST", "done calling db.Changes")

}

func TestCreateNeurgoCortex(t *testing.T) {
	game := &Game{}
	game.CreateNeurgoCortex()
	cortex := game.cortex
	assert.True(t, cortex != nil)
	assert.True(t, cortex.Sensors != nil)

	cortex.RenderSVGFile("out.svg")

}

func FakeGameDocument() (gameState []float64, possibleMoves []Move) {

	gameState = make([]float64, 32)

	possibleMove1 := Move{
		startLocation:   0,
		isCurrentlyKing: -1,
		endLocation:     1.0,
		willBecomeKing:  -0.5,
		captureValue:    1,
	}

	possibleMove2 := Move{
		startLocation:   1,
		isCurrentlyKing: -0.5,
		endLocation:     0.0,
		willBecomeKing:  0.5,
		captureValue:    0,
	}

	possibleMoves = []Move{possibleMove1, possibleMove2}
	return

}

func TestChooseBestMove(t *testing.T) {

	ng.SeedRandom()
	logg.LogKeys["MAIN"] = true

	game := &Game{}
	game.CreateNeurgoCortex()
	cortex := game.cortex
	cortex.Run()

	gameState, possibleMoves := FakeGameDocument()
	bestMove := game.ChooseBestMove(gameState, possibleMoves)

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
