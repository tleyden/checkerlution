package checkerlution

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/logg"
	"github.com/nu7hatch/gouuid"
	"github.com/tleyden/dsallings-couch-go"
	ng "github.com/tleyden/neurgo"
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
	cortex               *ng.Cortex
	currentGameState     GameStateVector
	gameState            GameState
	currentPossibleMove  ValidMoveCortexInput
	latestActuatorOutput []float64
	ourTeamId            int
	db                   couch.Database
	user                 User
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

// - make sure one of the changes is a game, if not, ignore it
// - get the latest game document
// - if it's not our turn, do nothing
// - if it is our turn
//   - parse out the required data structures needed to pass to cortex
//   - call cortex to calculate next move
//   - make next move by inserting a new revision of votes doc
func (game *Game) handleChanges(changes Changes) {
	logg.LogTo("DEBUG", "handleChanges called with %v", changes)
	gameDocChanged := game.checkGameDocInChanges(changes)
	if gameDocChanged {
		gameState, err := game.fetchLatestGameState()
		game.updateUserGameNumber(gameState)
		game.gameState = gameState
		if err != nil {
			logg.LogError(err)
			return
		}
		logg.LogTo("DEBUG", "gameState: %v", gameState)
		if isOurTurn := game.isOurTurn(gameState); !isOurTurn {
			logg.LogTo("DEBUG", "It's not our turn, ignoring changes")
			return
		}

		gameStateVector := game.extractGameStateVector(gameState)

		logg.LogTo("DEBUG", "gameStateVector: %v", gameStateVector)

		possibleMoves := game.extractPossibleMoves(gameState)

		if len(possibleMoves) == 0 {
			logg.LogTo("MAIN", "No possibleMoves, ignoring changes")
			return
		}

		logg.LogTo("DEBUG", "possibleMoves: %v", possibleMoves)

		bestMove := game.ChooseBestMove(gameStateVector, possibleMoves)

		logg.LogTo("DEBUG", "bestMove: %v", bestMove)

		game.PostChosenMove(bestMove)

	}

}

func (game Game) extractPossibleMoves(gameState GameState) []ValidMoveCortexInput {

	moves := make([]ValidMoveCortexInput, 0)

	ourTeam := gameState.Teams[game.ourTeamId]

	for pieceIndex, piece := range ourTeam.Pieces {
		piece.PieceId = pieceIndex
		for _, validMove := range piece.ValidMoves {
			moveInput := NewValidMoveCortexInput(validMove, piece)
			moves = append(moves, moveInput)
		}
	}

	return moves
}

func (game Game) opponentTeamId() int {
	switch game.ourTeamId {
	case RED_TEAM:
		return BLUE_TEAM
	default:
		return RED_TEAM
	}
}

func (game Game) extractGameStateVector(gameState GameState) GameStateVector {
	gameStateVector := NewGameStateVector()
	gameStateVector.loadFromGameState(gameState, game.ourTeamId)
	return gameStateVector
}

func (game Game) isOurTurn(gameState GameState) bool {
	return gameState.ActiveTeam == game.ourTeamId
}

func (game Game) checkGameDocInChanges(changes Changes) bool {
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

func (game *Game) InitGame() {
	game.CreateNeurgoCortex()
	cortex := game.cortex
	cortex.Run()
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

func (game *Game) ChooseBestMove(gameStateVector GameStateVector, possibleMoves []ValidMoveCortexInput) (bestMove ValidMoveCortexInput) {

	// Todo: the code below is an implementation of a single MoveChooser
	// but an interface should be designed so this is pluggable

	game.currentGameState = gameStateVector
	logg.LogTo("MAIN", "gameStateVector: %v", gameStateVector)

	var bestMoveRating []float64
	bestMoveRating = []float64{-1000000000}

	for _, move := range possibleMoves {

		logg.LogTo("MAIN", "feed possible move to cortex: %v", move)

		// present it to the neural net
		game.currentPossibleMove = move
		game.cortex.SyncSensors()
		game.cortex.SyncActuators()

		logg.LogTo("MAIN", "done sync'ing actuators")

		logg.LogTo("MAIN", "actuator output %v bestMoveRating: %v", game.latestActuatorOutput[0], bestMoveRating[0])
		if game.latestActuatorOutput[0] > bestMoveRating[0] {
			logg.LogTo("MAIN", "actuator output > bestMoveRating")
			bestMove = move
			bestMoveRating[0] = game.latestActuatorOutput[0]
		} else {
			logg.LogTo("MAIN", "actuator output < bestMoveRating, ignoring")
		}

	}
	return

}

func (game *Game) PostChosenMove(move ValidMoveCortexInput) {

	logg.LogTo("MAIN", "post chosen move: %v", move.validMove)

	preMoveSleepSeconds := game.calculatePreMoveSleepSeconds()

	logg.LogTo("MAIN", "sleep %v (s) before posting move", preMoveSleepSeconds)

	time.Sleep(time.Second * time.Duration(preMoveSleepSeconds))

	if len(move.validMove.Locations) == 0 {
		logg.LogTo("MAIN", "invalid move, ignoring: %v", move.validMove)
	}

	u4, err := uuid.NewV4()
	if err != nil {
		logg.LogPanic("Error generating uuid", err)
	}

	votes := &OutgoingVotes{}
	votes.Id = fmt.Sprintf("vote:%s", u4)
	votes.Turn = game.gameState.Turn
	votes.PieceId = move.validMove.PieceId
	votes.TeamId = game.ourTeamId
	votes.GameId = game.gameState.Number

	// TODO: this is actually a bug, because if there is a
	// double jump it will only send the first jump move
	endLocation := move.validMove.Locations[0]
	locations := []int{move.validMove.StartLocation, endLocation}
	votes.Locations = locations

	newId, newRevision, err := game.db.Insert(votes)

	logg.LogTo("MAIN", "newId: %v, newRevision: %v err: %v", newId, newRevision, err)
	if err != nil {
		logg.LogError(err)
		return
	}

}

func (game *Game) CreateNeurgoCortex() {

	nodeId := ng.NewCortexId("cortex")
	game.cortex = &ng.Cortex{
		NodeId: nodeId,
	}
	game.CreateSensors()
	game.CreateActuator()
	game.CreateNeuron()
	game.ConnectNodes()
}

func (game *Game) ConnectNodes() {

	cortex := game.cortex

	cortex.Init()

	// connect sensors -> neuron(s)
	for _, sensor := range cortex.Sensors {
		for _, neuron := range cortex.Neurons {
			sensor.ConnectOutbound(neuron)
			weights := ng.RandomWeights(sensor.VectorLength)
			neuron.ConnectInboundWeighted(sensor, weights)
		}
	}

	// connect neuron to actuator
	for _, neuron := range cortex.Neurons {
		for _, actuator := range cortex.Actuators {
			neuron.ConnectOutbound(actuator)
			actuator.ConnectInbound(neuron)
		}
	}

}

func (game *Game) CreateNeuron() {
	neuron := &ng.Neuron{
		ActivationFunction: ng.EncodableSigmoid(),
		NodeId:             ng.NewNeuronId("Neuron", 0.25),
		Bias:               ng.RandomBias(),
	}
	game.cortex.SetNeurons([]*ng.Neuron{neuron})
}

func (game *Game) CreateActuator() {

	actuatorNodeId := ng.NewActuatorId("Actuator", 0.5)
	actuatorFunc := func(outputs []float64) {
		logg.LogTo("MAIN", "actuator func called with: %v", outputs)
		game.latestActuatorOutput = outputs
		game.cortex.SyncChan <- actuatorNodeId // TODO: this should be in actuator itself, not in this function
	}
	actuator := &ng.Actuator{
		NodeId:           actuatorNodeId,
		VectorLength:     1,
		ActuatorFunction: actuatorFunc,
	}
	game.cortex.SetActuators([]*ng.Actuator{actuator})

}

func (game *Game) CreateSensors() {

	sensorLayer := 0.0

	sensorFuncGameState := func(syncCounter int) []float64 {
		logg.LogTo("MAIN", "sensor func game state called")
		return game.currentGameState
	}
	sensorGameStateNodeId := ng.NewSensorId("SensorGameState", sensorLayer)
	sensorGameState := &ng.Sensor{
		NodeId:         sensorGameStateNodeId,
		VectorLength:   32,
		SensorFunction: sensorFuncGameState,
	}

	sensorFuncPossibleMove := func(syncCounter int) []float64 {
		logg.LogTo("MAIN", "sensor func possible move called")
		return game.currentPossibleMove.VectorRepresentation()
	}
	sensorPossibleMoveNodeId := ng.NewSensorId("SensorPossibleMove", sensorLayer)
	sensorPossibleMove := &ng.Sensor{
		NodeId:         sensorPossibleMoveNodeId,
		VectorLength:   5, // start_location, is_king, final_location, will_be_king, amt_would_capture
		SensorFunction: sensorFuncPossibleMove,
	}
	game.cortex.SetSensors([]*ng.Sensor{sensorGameState, sensorPossibleMove})

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

	return ng.RandomInRange(minSleep, float64(game.gameState.MoveInterval))

}
