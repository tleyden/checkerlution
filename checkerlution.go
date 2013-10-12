package checkerlution

import (
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

const (
	RUNNING_MODE = iota
	TRAINING_MODE
)

type Checkerlution struct {
	ourTeamId            int
	cortex               *ng.Cortex
	currentGameState     GameStateVector
	currentPossibleMove  ValidMoveCortexInput
	latestActuatorOutput []float64
	mode                 int
	latestFitnessScore   float64
}

func (c *Checkerlution) Start(ourTeamId int) {
	c.ourTeamId = ourTeamId
	c.CreateNeurgoCortex()
	cortex := c.cortex
	cortex.Run()

}

func (c *Checkerlution) StartWithCortex(cortex *ng.Cortex, ourTeamId int) {
	c.ourTeamId = ourTeamId
	c.cortex = cortex
	cortex.Run()
}

func (c *Checkerlution) GameFinished(gameState cbot.GameState) (shouldQuit bool) {
	logg.LogTo("MAIN", "GameFinished")
	switch c.mode {
	case TRAINING_MODE:
		logg.LogTo("MAIN", "in TRAINING_MODE, so quitting")
		shouldQuit = true
		c.latestFitnessScore = c.calculateFitness(gameState)
	case RUNNING_MODE:
		logg.LogTo("MAIN", "in RUNNIG_MODE, so not quitting")
		shouldQuit = false
	}
	logg.LogTo("MAIN", "shouldQuit: %v", shouldQuit)
	return
}

func (c *Checkerlution) Think(gameState cbot.GameState) (bestMove cbot.ValidMove, ok bool) {
	ok = true
	gameStateVector := c.extractGameStateVector(gameState)
	possibleMoves := c.extractPossibleMoves(gameState)
	if len(possibleMoves) == 0 {
		logg.LogTo("MAIN", "No possibleMoves, ignoring changes")
		ok = false
		return
	}
	bestMoveCortex := c.chooseBestMove(gameStateVector, possibleMoves)
	bestMove = bestMoveCortex.validMove
	return
}

func (c Checkerlution) Cortex() *ng.Cortex {
	return c.cortex
}

func (c *Checkerlution) SetMode(mode int) {
	c.mode = mode
}

func (c Checkerlution) calculateFitness(gameState cbot.GameState) (fitness float64) {
	weWon := (gameState.WinningTeam == c.ourTeamId)
	switch weWon {
	case true:
		logg.LogTo("MAIN", "calculateFitness based on winning")
		// fitness will be a positive number
		// the least amount of moves we made, the higher the fitness
		fitness = 200
		fitness -= float64(gameState.Turn)
		if fitness < 1 {
			fitness = 1 // lowest possible fitness when winning
		}
	case false:
		logg.LogTo("MAIN", "calculateFitness based on losing")
		// fitness will be a negative number
		// the least amount of moves we made, the lower (more negative)
		// the fitness, because we didn't put up much of a fight
		fitness = -200
		fitness += float64(gameState.Turn)
		if fitness > -1 {
			fitness = -1 // highest possible fitness when losing
		}
	}

	logg.LogTo("MAIN", "calculateFitness returning: %v", fitness)
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
	logg.LogTo("MAIN", "gameStateVector: %v", gameStateVector)

	var bestMoveRating []float64
	bestMoveRating = []float64{-1000000000}

	for _, move := range possibleMoves {

		logg.LogTo("MAIN", "feed possible move to cortex: %v", move)

		// present it to the neural net
		c.currentPossibleMove = move
		c.cortex.SyncSensors()
		c.cortex.SyncActuators()

		logg.LogTo("MAIN", "done sync'ing actuators")

		logg.LogTo("MAIN", "actuator output %v bestMoveRating: %v", c.latestActuatorOutput[0], bestMoveRating[0])
		if c.latestActuatorOutput[0] > bestMoveRating[0] {
			logg.LogTo("MAIN", "actuator output > bestMoveRating")
			bestMove = move
			bestMoveRating[0] = c.latestActuatorOutput[0]
		} else {
			logg.LogTo("MAIN", "actuator output < bestMoveRating, ignoring")
		}

	}
	return

}

func (c *Checkerlution) CreateNeurgoCortex() {

	nodeId := ng.NewCortexId("cortex")
	c.cortex = &ng.Cortex{
		NodeId: nodeId,
	}
	c.CreateSensors()
	c.CreateActuator()
	c.CreateNeuron()
	c.ConnectNodes()
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
	actuatorFunc := func(outputs []float64) {
		logg.LogTo("MAIN", "actuator func called with: %v", outputs)
		c.latestActuatorOutput = outputs
	}
	actuator := &ng.Actuator{
		NodeId:           actuatorNodeId,
		VectorLength:     1,
		ActuatorFunction: actuatorFunc,
	}
	c.cortex.SetActuators([]*ng.Actuator{actuator})

}

func (c *Checkerlution) CreateSensors() {

	sensorLayer := 0.0

	sensorFuncGameState := func(syncCounter int) []float64 {
		logg.LogTo("MAIN", "sensor func game state called")
		return c.currentGameState
	}
	sensorGameStateNodeId := ng.NewSensorId("SensorGameState", sensorLayer)
	sensorGameState := &ng.Sensor{
		NodeId:         sensorGameStateNodeId,
		VectorLength:   32,
		SensorFunction: sensorFuncGameState,
	}

	sensorFuncPossibleMove := func(syncCounter int) []float64 {
		logg.LogTo("MAIN", "sensor func possible move called")
		return c.currentPossibleMove.VectorRepresentation()
	}
	sensorPossibleMoveNodeId := ng.NewSensorId("SensorPossibleMove", sensorLayer)
	sensorPossibleMove := &ng.Sensor{
		NodeId:         sensorPossibleMoveNodeId,
		VectorLength:   5, // start_location, is_king, final_location, will_be_king, amt_would_capture
		SensorFunction: sensorFuncPossibleMove,
	}
	c.cortex.SetSensors([]*ng.Sensor{sensorGameState, sensorPossibleMove})

}

func (c Checkerlution) Stop() {

}
