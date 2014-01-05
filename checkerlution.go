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

func (c *Checkerlution) CreateNeurgoCortex() {

	uuid := ng.NewUuid()
	cortexUuid := fmt.Sprintf("cortex-%s", uuid)
	nodeId := ng.NewCortexId(cortexUuid)

	c.cortex = &ng.Cortex{
		NodeId: nodeId,
	}

	c.cortex.Init()

	c.CreateSensors()

	outputNeuron := c.CreateOutputNeuron()
	layer1Neurons := c.CreateHiddenLayer1Neurons(outputNeuron)
	layer2Neurons := c.CreateHiddenLayer2Neurons(layer1Neurons, outputNeuron)

	// combine all into single slice and add neurons to cortex
	neurons := []*ng.Neuron{}
	neurons = append(neurons, layer1Neurons...)
	neurons = append(neurons, layer2Neurons...)
	neurons = append(neurons, outputNeuron)
	c.cortex.SetNeurons(neurons)

	actuator := c.CreateActuator()

	// workaround for error
	// Cannot make outbound connection, dataChan == nil [recovered]
	c.cortex.Init()

	outputNeuron.ConnectOutbound(actuator)
	actuator.ConnectInbound(outputNeuron)

	c.cortex.MarshalJSONToFile("checkerlution_cortex.json")

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

func (c *Checkerlution) CreateHiddenLayer1Neurons(outputNeuron *ng.Neuron) []*ng.Neuron {

	cortex := c.cortex
	neurons := []*ng.Neuron{}
	layerIndex := 0.25

	for i := 0; i < 40; i++ {
		name := fmt.Sprintf("hidden-layer-1-n-%d", i)
		neuron := &ng.Neuron{
			ActivationFunction: ng.EncodableTanh(),
			NodeId:             ng.NewNeuronId(name, layerIndex),
			Bias:               ng.RandomBias(),
		}

		// Workaround for error:
		// Cannot make outbound connection, dataChan == nil [recovered]
		// The best fix is to just load nn from json
		neuron.Init()

		for _, sensor := range cortex.Sensors {
			sensor.ConnectOutbound(neuron)
			weights := ng.RandomWeights(sensor.VectorLength)
			neuron.ConnectInboundWeighted(sensor, weights)
		}

		// connect directly to output neuron
		neuron.ConnectOutbound(outputNeuron)
		weights := ng.RandomWeights(1)
		outputNeuron.ConnectInboundWeighted(neuron, weights)

		neurons = append(neurons, neuron)

	}
	return neurons

}

func (c *Checkerlution) CreateHiddenLayer2Neurons(layer1Neurons []*ng.Neuron, outputNeuron *ng.Neuron) []*ng.Neuron {

	neurons := []*ng.Neuron{}
	layerIndex := 0.35

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("hidden-layer-1-n-%d", i)
		neuron := &ng.Neuron{
			ActivationFunction: ng.EncodableTanh(),
			NodeId:             ng.NewNeuronId(name, layerIndex),
			Bias:               ng.RandomBias(),
		}

		// Workaround for error:
		// Cannot make outbound connection, dataChan == nil [recovered]
		// The best fix is to just load nn from json
		neuron.Init()

		for _, layer1Neuron := range layer1Neurons {
			layer1Neuron.ConnectOutbound(neuron)
			weights := ng.RandomWeights(1)
			neuron.ConnectInboundWeighted(layer1Neuron, weights)
		}

		// connect directly to output neuron
		neuron.ConnectOutbound(outputNeuron)
		weights := ng.RandomWeights(1)
		outputNeuron.ConnectInboundWeighted(neuron, weights)

		neurons = append(neurons, neuron)

	}
	return neurons

}

func (c *Checkerlution) CreateOutputNeuron() *ng.Neuron {

	layerIndex := 0.45
	neuron := &ng.Neuron{
		ActivationFunction: ng.EncodableTanh(),
		NodeId:             ng.NewNeuronId("OutputNeuron", layerIndex),
		Bias:               ng.RandomBias(),
	}

	// Workaround for error:
	// Cannot make outbound connection, dataChan == nil [recovered]
	// The best fix is to just load nn from json
	neuron.Init()

	return neuron

}

func (c *Checkerlution) CreateActuator() *ng.Actuator {

	actuatorNodeId := ng.NewActuatorId("Actuator", 0.5)
	actuator := &ng.Actuator{
		NodeId:           actuatorNodeId,
		VectorLength:     1,
		ActuatorFunction: c.actuatorFunc(),
	}
	c.cortex.SetActuators([]*ng.Actuator{actuator})
	return actuator

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
	depth := 5 // TODO: crank this up higher
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
