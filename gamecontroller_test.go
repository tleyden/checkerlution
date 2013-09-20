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

func TestIsOurTurn(t *testing.T) {

	jsonString := `{"_id":"game:checkers","_rev":"3773-aa8a4c5a30b49e1eec65dff6df05561f","activeTeam":0,"channels":["game"],"moveDeadline":"2013-09-20T21:13:35Z","moveInterval":30,"moves":[{"game":153563,"locations":[10,14],"piece":9,"team":0,"turn":1},{"game":153563,"locations":[23,19],"piece":2,"team":1,"turn":2}],"number":153563,"startTime":"2013-09-20T17:11:53Z","teams":[{"participantCount":1,"pieces":[{"location":1},{"location":2},{"location":3},{"location":4},{"location":5},{"location":6,"validMoves":[{"captures":[],"king":false,"locations":[10]}]},{"location":7,"validMoves":[{"captures":[],"king":false,"locations":[10]}]},{"location":8},{"location":9,"validMoves":[{"captures":[],"king":false,"locations":[13]}]},{"location":14,"validMoves":[{"captures":[],"king":false,"locations":[17]},{"captures":[],"king":false,"locations":[18]}]},{"location":11,"validMoves":[{"captures":[],"king":false,"locations":[15]},{"captures":[],"king":false,"locations":[16]}]},{"location":12,"validMoves":[{"captures":[],"king":false,"locations":[16]}]}]},{"participantCount":0,"pieces":[{"location":21},{"location":22},{"location":19},{"location":24},{"location":25},{"location":26},{"location":27},{"location":28},{"location":29},{"location":30},{"location":31},{"location":32}]}],"turn":3,"votesDoc":"votes:checkers"}`

	jsonBytes := []byte(jsonString)
	gameDoc := new(Document)
	err := json.Unmarshal(jsonBytes, gameDoc)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{ourTeamId: 0}
	game.InitGame()
	result := game.isOurTurn(*gameDoc)
	assert.True(t, result)

	game.ourTeamId = 1
	assert.False(t, game.isOurTurn(*gameDoc))

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

func DISTestFetchLatestGameDocument(t *testing.T) {
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
