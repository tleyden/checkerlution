package checkerlution

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

type Checkerlution struct {
	game                 Game
	cortex               *ng.Cortex
	currentGameState     GameStateVector
	currentPossibleMove  ValidMoveCortexInput
	latestActuatorOutput []float64
}

func (c *Checkerlution) Start(game Game) {
	c.game = game
	c.CreateNeurgoCortex()
	cortex := c.cortex
	cortex.Run()

}

func (c *Checkerlution) Think(gameState GameState) (bestMove ValidMove) {

	gameStateVector := c.extractGameStateVector(gameState)
	possibleMoves := c.extractPossibleMoves(gameState)
	if len(possibleMoves) == 0 {
		logg.LogTo("MAIN", "No possibleMoves, ignoring changes")
		return
	}
	bestMoveCortex := c.chooseBestMove(gameStateVector, possibleMoves)
	bestMove = bestMoveCortex.validMove
	return
}

func (c Checkerlution) extractGameStateVector(gameState GameState) GameStateVector {
	gameStateVector := NewGameStateVector()
	gameStateVector.loadFromGameState(gameState, c.game.ourTeamId)
	return gameStateVector
}

func (c Checkerlution) extractPossibleMoves(gameState GameState) []ValidMoveCortexInput {

	moves := make([]ValidMoveCortexInput, 0)

	ourTeam := gameState.Teams[c.game.ourTeamId]

	for pieceIndex, piece := range ourTeam.Pieces {
		piece.PieceId = pieceIndex
		for _, validMove := range piece.ValidMoves {
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
		ActivationFunction: ng.EncodableSigmoid(),
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
		c.cortex.SyncChan <- actuatorNodeId // TODO: this should be in actuator itself, not in this function
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
