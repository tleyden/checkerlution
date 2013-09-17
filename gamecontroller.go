package checkerlution

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

type Game struct {
	cortex               *ng.Cortex
	currentGameState     []float64
	currentPossibleMove  []float64
	latestActuatorOutput []float64
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
		game.latestActuatorOutput = outputs
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
		return game.currentGameState
	}
	sensorGameStateNodeId := ng.NewSensorId("SensorGameState", sensorLayer)
	sensorGameState := &ng.Sensor{
		NodeId:         sensorGameStateNodeId,
		VectorLength:   32,
		SensorFunction: sensorFuncGameState,
	}

	sensorFuncPossibleMove := func(syncCounter int) []float64 {
		return game.currentPossibleMove
	}
	sensorPossibleMoveNodeId := ng.NewSensorId("SensorPossibleMove", sensorLayer)
	sensorPossibleMove := &ng.Sensor{
		NodeId:         sensorPossibleMoveNodeId,
		VectorLength:   5, // start_location, is_king, final_location, will_be_king, amt_would_capture
		SensorFunction: sensorFuncPossibleMove,
	}
	game.cortex.SetSensors([]*ng.Sensor{sensorGameState, sensorPossibleMove})

}

func (game *Game) GameLoop() {

	// get a neurgo network
	game.CreateNeurgoCortex()
	logg.LogTo("DEBUG", "game: %v", game)

	for {
		// read a new Game document via http

		// extract game state and list of available moves

		// for each available move:

		// present it to the neural net by saving in game
		// and calling cortex.SyncSensors

		// read response from actuator and store

		// post the chosen move onto a channel

		break

	}

}
