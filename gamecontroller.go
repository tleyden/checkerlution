package checkerlution

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

type Game struct {
	cortex               *ng.Cortex
	currentGameState     []float64
	currentPossibleMove  Move
	latestActuatorOutput []float64
}

func (game *Game) GameLoop() {

	client := Client{}

	// get a neurgo network
	game.CreateNeurgoCortex()
	cortex := game.cortex

	cortex.Run()

	for {

		// fetch game state and list of available moves from game server
		gameState, possibleMoves := client.FetchNewGameDocument()

		bestMove := game.ChooseBestMove(gameState, possibleMoves)

		game.PostChosenMove(bestMove)

		// when do we break out of the loop??

	}

	game.cortex.Shutdown()

}

func (game *Game) ChooseBestMove(gameState []float64, possibleMoves []Move) (bestMove Move) {

	game.currentGameState = gameState
	logg.LogTo("MAIN", "gameState: %v", gameState)

	var bestMoveRating []float64
	bestMoveRating = []float64{-1000000000}

	for _, move := range possibleMoves {

		logg.LogTo("MAIN", "possible move: %v", move)

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

func (game *Game) PostChosenMove(move Move) {
	logg.LogTo("MAIN", "chosen move: %v", move)
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
