package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"log"
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

func TestCheckGameDocInChanges(t *testing.T) {

	jsonString := `{"results":[{"seq":"*:3408","id":"user:6213C1A1-4E5F-429E-91C9-CDC2BF1537C3","changes":[{"rev":"3-783b9cda9b7b9e6faac2d8bda9e16535"}]},{"seq":"*:3409","id":"vote:6213C1A1-4E5F-429E-91C9-CDC2BF1537C3","changes":[{"rev":"1-393aaf8f37404c4a0159d9ec8dc1e0ee"}]},{"seq":"*:3440","id":"votes:checkers","changes":[{"rev":"16-ebaa86d97e63940fddfdbd11a219e9e6"}]},{"seq":"*:3641","id":"game:checkers","changes":[{"rev":"3586-09a232e6b524940185b0b268483981ea"}]}],"last_seq":"*:3641"}`
	jsonBytes := []byte(jsonString)
	changes := new(Changes)
	err := json.Unmarshal(jsonBytes, changes)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{}
	game.InitGame()
	result := game.checkGameDocInChanges(*changes)
	assert.True(t, result)

}

func TestFetchLatestGameDocument(t *testing.T) {
	game := &Game{}
	game.InitGame()
	gameDoc, err := game.fetchLatestGameDoc()
	if err != nil {
		logg.Log("gameDoc: %v.  err: %v", gameDoc, err)
		panic("err")
	}
	logg.Log("gameDoc: %v.  err: %v", gameDoc, err)
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

func TestGameLoop(t *testing.T) {
	ng.SeedRandom()
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	game := &Game{}
	game.GameLoop()

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
