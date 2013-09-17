package checkerlution

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

type Game struct {
	cortex *ng.Cortex
}

func (game *Game) CreateNeurgoCortex() {

	nodeId := ng.NewCortexId("cortex")
	game.cortex = &ng.Cortex{
		NodeId: nodeId,
	}
	game.CreateSensors()
}

func (game *Game) CreateSensors() {

	sensorFunc := func(syncCounter int) []float64 {
		return []float64{0}
	}

	sensorLayer := 0.0
	sensorGameStateNodeId := ng.NewSensorId("SensorGameState", sensorLayer)
	sensorGameState := &ng.Sensor{
		NodeId:         sensorGameStateNodeId,
		VectorLength:   32,
		SensorFunction: sensorFunc,
	}

	sensorPossibleMoveNodeId := ng.NewSensorId("SensorPossibleMove", sensorLayer)
	sensorPossibleMove := &ng.Sensor{
		NodeId:         sensorPossibleMoveNodeId,
		VectorLength:   5, // start_location, is_king, final_location, will_be_king, amt_would_capture
		SensorFunction: sensorFunc,
	}
	game.cortex.SetSensors([]*ng.Sensor{sensorGameState, sensorPossibleMove})

}

func (game *Game) GameLoop() {

	// get a neurgo network
	game.CreateNeurgoCortex()
	logg.LogTo("DEBUG", "game: %v", game)

	for {
		// read a new Game document (from a channel, I guess)

		// extract game state and list of available moves

		// for each available move:

		// present it to the neural net

		// read response from actuator and store

		// post the chosen move onto a channel

		break

	}

}
