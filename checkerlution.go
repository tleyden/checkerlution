package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	core "github.com/tleyden/checkers-core"
	ng "github.com/tleyden/neurgo"
)

type OperationMode int

const (
	RUNNING_MODE = iota
	TRAINING_MODE
)

type Checkerlution struct {
	ourTeamId            cbot.TeamType
	cortex               *ng.Cortex
	currentGameState     GameStateVector
	currentPossibleMove  ValidMoveCortexInput
	latestActuatorOutput []float64
	mode                 OperationMode
	latestFitnessScore   float64
}

func (c *Checkerlution) Start(ourTeamId cbot.TeamType) {
	c.ourTeamId = ourTeamId
	c.CreateNeurgoCortex()
	cortex := c.cortex
	cortex.Run()

}

func (c *Checkerlution) StartWithCortex(cortex *ng.Cortex, ourTeamId cbot.TeamType) {
	c.ourTeamId = ourTeamId
	c.setSensorActuatorFunctions(cortex)
	c.cortex = cortex
	cortex.Run()
}

func (c *Checkerlution) ThinkOLD(gameState cbot.GameState) (bestMove cbot.ValidMove, ok bool) {
	ok = true
	gameStateVector := c.extractGameStateVector(gameState)
	possibleMoves := c.extractPossibleMoves(gameState)
	if len(possibleMoves) == 0 {
		logg.LogTo("CHECKERLUTION", "No possibleMoves, ignoring changes")
		ok = false
		return
	}
	bestMoveCortex := c.chooseBestMove(gameStateVector, possibleMoves)
	bestMove = bestMoveCortex.validMove
	return
}

func (c *Checkerlution) Think(gameState cbot.GameState) (bestMove cbot.ValidMove, ok bool) {

	ok = true
	ourTeam := gameState.Teams[c.ourTeamId]
	allValidMoves := ourTeam.AllValidMoves()
	if len(allValidMoves) > 0 {

		// convert into core.board representation
		board := gameState.Export()
		logg.LogTo("DEBUG", "Before move %v", board.CompactString(true))

		// generate best move (will be a core.move) -- initially, pick random
		move := c.generateBestMove(board)

		// search allValidMoves to find corresponding valid move
		found, bestValidMoveIndex := cbot.CorrespondingValidMoveIndex(move, allValidMoves)

		if !found {
			msg := "Could not find corresponding valid move: %v in %v"
			logg.LogPanic(msg, move, allValidMoves)
		} else {
			bestMove = allValidMoves[bestValidMoveIndex]
		}

		// this is just for debugging purposes
		player := cbot.GetCorePlayer(c.ourTeamId)
		boardPostMove := board.ApplyMove(player, move)
		logg.LogTo("DEBUG", "After move %v", boardPostMove.CompactString(true))

		return

	} else {
		ok = false
	}

	return

}

func (c *Checkerlution) GameFinished(gameState cbot.GameState) (shouldQuit bool) {
	switch c.mode {
	case TRAINING_MODE:
		shouldQuit = true
		c.latestFitnessScore = c.calculateFitness(gameState)
	case RUNNING_MODE:
		shouldQuit = false
	}
	return
}

func (c Checkerlution) Cortex() *ng.Cortex {
	return c.cortex
}

func (c *Checkerlution) SetMode(mode OperationMode) {
	c.mode = mode
}

func (c Checkerlution) calculateFitness(gameState cbot.GameState) (fitness float64) {
	weWon := (gameState.WinningTeam == c.ourTeamId)
	switch weWon {
	case true:
		logg.LogTo("CHECKERLUTION", "calculateFitness based on winning.  Turn: %v", gameState.Turn)
		// fitness will be a positive number
		// the least amount of moves we made, the higher the fitness
		fitness = 200
		fitness -= float64(gameState.Turn)
		if fitness < 1 {
			fitness = 1 // lowest possible fitness when winning
		}
	case false:
		logg.LogTo("CHECKERLUTION", "calculateFitness based on losing.  Turn: %v", gameState.Turn)
		// fitness will be a negative number
		// the least amount of moves we made, the lower (more negative)
		// the fitness, because we didn't put up much of a fight
		fitness = -200
		fitness += float64(gameState.Turn)
		if fitness > -1 {
			fitness = -1 // highest possible fitness when losing
		}
	}

	logg.LogTo("CHECKERLUTION", "calculateFitness returning: %v", fitness)
	return
}

func (c Checkerlution) extractGameStateVector(gameState cbot.GameState) GameStateVector {
	gameStateVector := NewGameStateVector()
	gameStateVector.loadFromGameState(gameState, c.ourTeamId)
	return gameStateVector
}

func (c Checkerlution) extractPossibleMoves(gameState cbot.GameState) []ValidMoveCortexInput {

	moves := make([]ValidMoveCortexInput, 0)

	ourTeam := gameState.Teams[c.ourTeamId]

	for pieceIndex, piece := range ourTeam.Pieces {
		for _, validMove := range piece.ValidMoves {

			// enhance the validMove from some information
			// from the piece, because this will be required
			// later when translating this valid move into an
			// "outgoing move", eg, a move that can be posted
			// to server to cause it to make.  the outgoing move
			// is a pretty different format than the orig validMove
			validMove.StartLocation = piece.Location
			validMove.PieceId = pieceIndex

			moveInput := NewValidMoveCortexInput(validMove, piece)
			moves = append(moves, moveInput)
		}
	}

	return moves
}

func (c *Checkerlution) chooseBestMove(gameStateVector GameStateVector, possibleMoves []ValidMoveCortexInput) (bestMove ValidMoveCortexInput) {

	c.currentGameState = gameStateVector
	logg.LogTo("CHECKERLUTION", "gameStateVector: %v", gameStateVector)

	var bestMoveRating []float64
	bestMoveRating = []float64{-1000000000}

	for _, move := range possibleMoves {

		logg.LogTo("CHECKERLUTION", "feed possible move to cortex: %v", move)

		// present it to the neural net
		c.currentPossibleMove = move
		c.cortex.SyncSensors()
		c.cortex.SyncActuators()

		logg.LogTo("CHECKERLUTION", "actuator output %v bestMoveRating: %v", c.latestActuatorOutput[0], bestMoveRating[0])
		if c.latestActuatorOutput[0] > bestMoveRating[0] {
			logg.LogTo("CHECKERLUTION", "actuator output > bestMoveRating")
			bestMove = move
			bestMoveRating[0] = c.latestActuatorOutput[0]
		} else {
			logg.LogTo("CHECKERLUTION", "actuator output < bestMoveRating, ignoring")
		}

	}
	return

}

func (c *Checkerlution) CreateNeurgoCortex() {

	uuid := ng.NewUuid()
	cortexUuid := fmt.Sprintf("cortex-%s", uuid)
	nodeId := ng.NewCortexId(cortexUuid)

	c.cortex = &ng.Cortex{
		NodeId: nodeId,
	}
	c.CreateSensors()
	c.CreateActuator()
	c.CreateNeuron()
	c.ConnectNodes()
}

func (c *Checkerlution) LoadNeurgoCortex(filename string) {

	cortex, err := ng.NewCortexFromJSONFile(filename)
	if err != nil {
		logg.LogPanic("Error reading cortex from: %v.  Err: %v", filename, err)
	}

	c.setSensorActuatorFunctions(cortex)

	c.cortex = cortex
}

func (c *Checkerlution) setSensorActuatorFunctions(cortex *ng.Cortex) {

	sensor := cortex.FindSensor(ng.NewSensorId("SensorGameState", 0))
	sensor.SensorFunction = c.sensorFuncGameState()
	sensor = cortex.FindSensor(ng.NewSensorId("SensorPossibleMove", 0))
	sensor.SensorFunction = c.sensorFuncPossibleMove()

	actuator := cortex.FindActuator(ng.NewActuatorId("Actuator", 0))
	actuator.ActuatorFunction = c.actuatorFunc()

}

func (c *Checkerlution) ConnectNodes() {

	cortex := c.cortex

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

func (c *Checkerlution) CreateNeuron() {
	neuron := &ng.Neuron{
		ActivationFunction: ng.EncodableTanh(),
		NodeId:             ng.NewNeuronId("Neuron", 0.25),
		Bias:               ng.RandomBias(),
	}
	c.cortex.SetNeurons([]*ng.Neuron{neuron})
}

func (c *Checkerlution) CreateActuator() {

	actuatorNodeId := ng.NewActuatorId("Actuator", 0.5)
	actuator := &ng.Actuator{
		NodeId:           actuatorNodeId,
		VectorLength:     1,
		ActuatorFunction: c.actuatorFunc(),
	}
	c.cortex.SetActuators([]*ng.Actuator{actuator})

}

func (c *Checkerlution) actuatorFunc() ng.ActuatorFunction {
	return func(outputs []float64) {
		logg.LogTo("CHECKERLUTION", "actuator func called with: %v", outputs)
		c.latestActuatorOutput = outputs
	}
}

func (c *Checkerlution) CreateSensors() {

	sensorLayer := 0.0

	sensorGameStateNodeId := ng.NewSensorId("SensorGameState", sensorLayer)
	sensorGameState := &ng.Sensor{
		NodeId:         sensorGameStateNodeId,
		VectorLength:   32,
		SensorFunction: c.sensorFuncGameState(),
	}

	c.cortex.SetSensors([]*ng.Sensor{sensorGameState})

}

func (c *Checkerlution) sensorFuncGameState() ng.SensorFunction {
	return func(syncCounter int) []float64 {
		logg.LogTo("CHECKERLUTION", "sensor func game state called on thinker: %p, returning: %v", c, c.currentGameState)
		if len(c.currentGameState) == 0 {
			logg.LogPanic("sensor would return invalid gamestate")
		}
		return c.currentGameState
	}
}

func (c *Checkerlution) sensorFuncPossibleMove() ng.SensorFunction {
	return func(syncCounter int) []float64 {
		logg.LogTo("CHECKERLUTION", "sensor func possible move called")
		return c.currentPossibleMove.VectorRepresentation()
	}
}

func (c Checkerlution) Stop() {

}

func (c *Checkerlution) generateBestMove(board core.Board) core.Move {
	evalFunc := c.getEvaluationFunction()
	player := cbot.GetCorePlayer(c.ourTeamId)
	depth := 6 // TODO: crank this up higher
	bestMove, scorePostMove := board.Minimax(player, depth, evalFunc)
	logg.LogTo("DEBUG", "scorePostMove: %v", scorePostMove)
	return bestMove
}

func (c *Checkerlution) getEvaluationFunction() core.EvaluationFunction {

	evalFunc := func(currentPlayer core.Player, board core.Board) float64 {

		// convert the board into inputs for the neural net (32 elt vector)
		// taking into account whether this player is "us" or not
		gameStateVector := NewGameStateVector()
		gameStateVector.loadFromBoard(board, currentPlayer)

		// send input to the neural net
		logg.LogTo("CHECKERLUTION", "set currentGameState %v", gameStateVector)
		c.currentGameState = gameStateVector
		c.cortex.SyncSensors()
		c.cortex.SyncActuators()

		// get output
		logg.LogTo("CHECKERLUTION", "actuator output %v", c.latestActuatorOutput[0])

		// return output
		return c.latestActuatorOutput[0]

	}
	return evalFunc

}
