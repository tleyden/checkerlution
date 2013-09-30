package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"log"
	"testing"
)

func init() {

	logg.LogKeys["NODE_SEND"] = true
	logg.LogKeys["NODE_RECV"] = true
	logg.LogKeys["TEST"] = true
	logg.LogKeys["DEBUG"] = true

}

func TestIsOurTurn(t *testing.T) {

	jsonString := `{"_id":"game:checkers","_rev":"3773-aa8a4c5a30b49e1eec65dff6df05561f","activeTeam":0,"channels":["game"],"moveDeadline":"2013-09-20T21:13:35Z","moveInterval":30,"moves":[{"game":153563,"locations":[10,14],"piece":9,"team":0,"turn":1},{"game":153563,"locations":[23,19],"piece":2,"team":1,"turn":2}],"number":153563,"startTime":"2013-09-20T17:11:53Z","teams":[{"participantCount":1,"pieces":[{"location":1},{"location":2},{"location":3},{"location":4},{"location":5},{"location":6,"validMoves":[{"captures":[],"king":false,"locations":[10]}]},{"location":7,"validMoves":[{"captures":[],"king":false,"locations":[10]}]},{"location":8},{"location":9,"validMoves":[{"captures":[],"king":false,"locations":[13]}]},{"location":14,"validMoves":[{"captures":[],"king":false,"locations":[17]},{"captures":[],"king":false,"locations":[18]}]},{"location":11,"validMoves":[{"captures":[],"king":false,"locations":[15]},{"captures":[],"king":false,"locations":[16]}]},{"location":12,"validMoves":[{"captures":[],"king":false,"locations":[16]}]}]},{"participantCount":0,"pieces":[{"location":21},{"location":22},{"location":19},{"location":24},{"location":25},{"location":26},{"location":27},{"location":28},{"location":29},{"location":30},{"location":31},{"location":32}]}],"turn":3,"votesDoc":"votes:checkers"}`

	gameState := NewGameStateFromString(jsonString)

	game := &Game{ourTeamId: 0}
	result := game.isOurTurn(gameState)
	assert.True(t, result)

	game.ourTeamId = 1
	assert.False(t, game.isOurTurn(gameState))

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
	result := game.checkGameDocInChanges(*changes)
	assert.True(t, result)

}

func TestCalculatePreMoveSleepSeconds(t *testing.T) {
	game := &Game{}
	game.gameState.MoveInterval = 30
	preMoveSleepSeconds := game.calculatePreMoveSleepSeconds()
	assert.True(t, preMoveSleepSeconds > 0)
	assert.True(t, preMoveSleepSeconds <= 30)
}
