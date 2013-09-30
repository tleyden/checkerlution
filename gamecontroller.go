package checkerlution

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/nu7hatch/gouuid"
	"github.com/tleyden/dsallings-couch-go"
	"io"
	"strings"
	"time"
)

const (
	SERVER_URL   = "http://localhost:4984/checkers"
	GAME_DOC_ID  = "game:checkers"
	VOTES_DOC_ID = "votes:checkers"
	RED_TEAM     = 0
	BLUE_TEAM    = 1
)

type Game struct {
	thinker   Thinker
	gameState GameState
	ourTeamId int
	db        couch.Database
	user      User
}

type Changes map[string]interface{}

func NewGame(ourTeamId int) *Game {
	game := &Game{ourTeamId: ourTeamId}
	return game
}

// Follow the changes feed and on each change callback
// call game.handleChanges() which will drive the game
func (game *Game) GameLoop() {

	game.InitGame()

	curSinceValue := "0"

	handleChange := func(reader io.Reader) string {
		changes := decodeChanges(reader)
		game.handleChanges(changes)
		curSinceValue = calculateNextSinceValue(curSinceValue, changes)
		time.Sleep(time.Second * 5)
		return curSinceValue
	}

	options := Changes{"since": "0"}
	game.db.Changes(handleChange, options)

}

// - make sure one of the changes is a game, if not, ignore it
// - get the latest game document
// - if it's not our turn, do nothing
// - if it is our turn
//   - call thinker to calculate next move
//   - make next move by inserting a new revision of votes doc
func (game *Game) handleChanges(changes Changes) {
	gameDocChanged := game.hasGameDocChanged(changes)
	if gameDocChanged {
		gameState, err := game.fetchLatestGameState()
		if err != nil {
			logg.LogError(err)
			return
		}
		game.updateUserGameNumber(gameState)
		game.gameState = gameState
		if isOurTurn := game.isOurTurn(gameState); !isOurTurn {
			logg.LogTo("DEBUG", "It's not our turn, ignoring changes")
			return
		}
		bestMove := game.thinker.Think(gameState)
		game.PostChosenMove(bestMove)

	}
}

func (game *Game) InitGame() {

	// game.thinker = &Checkerlution{}
	game.thinker = &RandomThinker{}

	game.thinker.Start(game.ourTeamId)
	game.InitDbConnection()
	game.CreateRemoteUser()
}

func (game *Game) CreateRemoteUser() {

	u4, err := uuid.NewV4()
	if err != nil {
		logg.LogPanic("Error generating uuid", err)
	}

	user := &User{
		Id:     fmt.Sprintf("user:%s", u4),
		TeamId: game.ourTeamId,
	}
	newId, newRevision, err := game.db.Insert(user)
	logg.LogTo("MAIN", "Inserted new user %v rev %v", newId, newRevision)

	user.Rev = newRevision
	game.user = *user

}

func (game *Game) InitDbConnection() {
	db, error := couch.Connect(SERVER_URL)
	if error != nil {
		logg.LogPanic("Error connecting to %v: %v", SERVER_URL, error)
	}
	game.db = db
}

func (game *Game) PostChosenMove(validMove ValidMove) {

	logg.LogTo("MAIN", "post chosen move: %v", validMove)

	preMoveSleepSeconds := game.calculatePreMoveSleepSeconds()

	logg.LogTo("MAIN", "sleep %v (s) before posting move", preMoveSleepSeconds)

	time.Sleep(time.Second * time.Duration(preMoveSleepSeconds))

	if len(validMove.Locations) == 0 {
		logg.LogTo("MAIN", "invalid move, ignoring: %v", validMove)
	}

	u4, err := uuid.NewV4()
	if err != nil {
		logg.LogPanic("Error generating uuid", err)
	}

	votes := &OutgoingVotes{}
	votes.Id = fmt.Sprintf("vote:%s", u4)
	votes.Turn = game.gameState.Turn
	votes.PieceId = validMove.PieceId
	votes.TeamId = game.ourTeamId
	votes.GameId = game.gameState.Number

	// TODO: this is actually a bug, because if there is a
	// double jump it will only send the first jump move
	endLocation := validMove.Locations[0]
	locations := []int{validMove.StartLocation, endLocation}
	votes.Locations = locations

	newId, newRevision, err := game.db.Insert(votes)

	logg.LogTo("MAIN", "newId: %v, newRevision: %v err: %v", newId, newRevision, err)
	if err != nil {
		logg.LogError(err)
		return
	}

}

func (game *Game) updateUserGameNumber(gameState GameState) {
	gameNumberChanged := (game.gameState.Number != gameState.Number)
	if gameNumberChanged {
		game.user.GameNumber = gameState.Number
		newRevision, err := game.db.Edit(game.user)
		if err != nil {
			logg.LogError(err)
			return
		}
		logg.LogTo("MAIN", "user update, rev: %v", newRevision)
	}

}

func (game Game) opponentTeamId() int {
	switch game.ourTeamId {
	case RED_TEAM:
		return BLUE_TEAM
	default:
		return RED_TEAM
	}
}

func (game Game) isOurTurn(gameState GameState) bool {
	return gameState.ActiveTeam == game.ourTeamId
}

func (game Game) hasGameDocChanged(changes Changes) bool {
	foundGameDoc := false
	changeResultsRaw := changes["results"]
	changeResults := changeResultsRaw.([]interface{})
	for _, changeResultRaw := range changeResults {
		changeResult := changeResultRaw.(map[string]interface{})
		docIdRaw := changeResult["id"]
		docId := docIdRaw.(string)
		if strings.Contains(docId, GAME_DOC_ID) {
			foundGameDoc = true
		}
	}
	return foundGameDoc
}

func (game Game) fetchLatestGameState() (gameState GameState, err error) {
	gameStateFetched := &GameState{}
	err = game.db.Retrieve(GAME_DOC_ID, gameStateFetched)
	if err == nil {
		gameState = *gameStateFetched
	}
	return
}

func decodeChanges(reader io.Reader) Changes {
	changes := make(Changes)
	decoder := json.NewDecoder(reader)
	decoder.Decode(&changes)
	return changes
}

func calculateNextSinceValue(curSinceValue string, changes Changes) string {
	lastSeq := changes["last_seq"]
	lastSeqAsString := lastSeq.(string)
	if lastSeq != nil && len(lastSeqAsString) > 0 {
		return lastSeqAsString
	}
	return curSinceValue
}

func (game *Game) calculatePreMoveSleepSeconds() float64 {

	// we don't want to make a move "too soon", so lets
	// cap the minimum amount we sleep at 10% of the move interval
	minSleep := float64(game.gameState.MoveInterval) * 0.10

	// likewise, don't want to cut it to close to the timeout
	maxSleep := float64(game.gameState.MoveInterval) * 0.90

	return randomInRange(minSleep, maxSleep)

}
