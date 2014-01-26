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

type CheckerlutionFlags struct {
	CheckersBotFlags cbot.CheckersBotFlags
	PopulationName   string
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
		logg.LogTo("DEBUG", "Before team %v move %v", c.ourTeamId.String(), board.CompactString(true))

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
		logg.LogTo("DEBUG", "After team %v move %v", c.ourTeamId.String(), boardPostMove.CompactString(true))

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

	switch gameState.WinningTeam {
	case c.ourTeamId:
		logg.LogTo("CHECKERLUTION", "calculateFitness based on winning.  Turn: %v", gameState.Turn)
		fitness = 1.0
	case c.ourTeamId.Opponent():
		logg.LogTo("CHECKERLUTION", "calculateFitness based on losing.  Turn: %v", gameState.Turn)
		fitness = -2.0
	default:
		logg.LogTo("CHECKERLUTION", "calculateFitness based on draw.  Turn: %v", gameState.Turn)
		fitness = 0.0
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

	actuator := cortex.FindActuator(ng.NewActuatorId("Actuator", 0))
	actuator.ActuatorFunction = c.actuatorFunc()

}

func (c *Checkerlution) CreateHiddenLayer1Neurons(outputNeuron *ng.Neuron) []*ng.Neuron {

	cortex := c.cortex
	neurons := []*ng.Neuron{}
	layerIndex := 0.25

	for i := 0; i < 15; i++ {
		name := fmt.Sprintf("hidden-layer-%f-n-%d", layerIndex, i)
		neuron := &ng.Neuron{
			ActivationFunction: ng.EncodableTanh(),
			NodeId:             ng.NewNeuronId(name, layerIndex),
			Bias:               ng.RandomBias(),
		}

		// Workaround for error:
		// Cannot make outbound connection, dataChan == nil [recovered]
		neuron.Init()

		for _, sensor := range cortex.Sensors {
			sensor.ConnectOutbound(neuron)
			weights := ng.RandomWeights(sensor.VectorLength)
			neuron.ConnectInboundWeighted(sensor, weights)
		}

		neurons = append(neurons, neuron)

	}
	return neurons

}

func (c *Checkerlution) CreateHiddenLayer2Neurons(layer1Neurons []*ng.Neuron, outputNeuron *ng.Neuron) []*ng.Neuron {

	neurons := []*ng.Neuron{}
	layerIndex := 0.35

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("hidden-layer-%f-n-%d", layerIndex, i)
		neuron := &ng.Neuron{
			ActivationFunction: ng.EncodableTanh(),
			NodeId:             ng.NewNeuronId(name, layerIndex),
			Bias:               ng.RandomBias(),
		}

		// Workaround for error:
		// Cannot make outbound connection, dataChan == nil [recovered]
		neuron.Init()

		for _, layer1Neuron := range layer1Neurons {
			layer1Neuron.ConnectOutbound(neuron)
			weights := ng.RandomWeights(1)
			neuron.ConnectInboundWeighted(layer1Neuron, weights)
		}

		// connect to output neuron
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

	// connect sensor directly to output neuron
	for _, sensor := range c.cortex.Sensors {
		sensor.ConnectOutbound(neuron)
		weights := ng.RandomWeights(sensor.VectorLength)
		neuron.ConnectInboundWeighted(sensor, weights)
	}

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
		if len(c.currentGameState) == 0 {
			logg.LogPanic("sensor would return invalid gamestate")
		}
		return c.currentGameState
	}
}

func (c *Checkerlution) sensorFuncPossibleMove() ng.SensorFunction {
	return func(syncCounter int) []float64 {
		return c.currentPossibleMove.VectorRepresentation()
	}
}

func (c Checkerlution) Stop() {

}

func (c *Checkerlution) generateBestMove(board core.Board) core.Move {

	counter := 0
	evalFunc := c.getEvaluationFunction(&counter)
	player := cbot.GetCorePlayer(c.ourTeamId)

	// with depth = 5, not working too well on first move.  when
	// there are only a few pieces on the board it seems to work,
	// but with full board .. taking a long time.

	depth := 4 // TODO: crank this up higher
	bestMove, scorePostMove := board.Minimax(player, depth, evalFunc)
	logg.LogTo("DEBUG", "scorePostMove: %v.  boards eval'd: %v", scorePostMove, counter)
	return bestMove
}

func (c *Checkerlution) getEvaluationFunction(counter *int) core.EvaluationFunction {

	evalFunc := func(currentPlayer core.Player, board core.Board) float64 {

		*counter += 1

		// convert the board into inputs for the neural net (32 elt vector)
		// taking into account whether this player is "us" or not
		gameStateVector := NewGameStateVector()
		gameStateVector.loadFromBoard(board, currentPlayer)

		// send input to the neural net
		c.currentGameState = gameStateVector
		c.cortex.SyncSensors()
		c.cortex.SyncActuators()

		// return output
		return c.latestActuatorOutput[0]

	}
	return evalFunc

}
