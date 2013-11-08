package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["NEURGO"] = true
	logg.LogKeys["SENSOR_SYNC"] = true
	logg.LogKeys["ACTUATOR_SYNC"] = true
	logg.LogKeys["NODE_PRE_SEND"] = true
	logg.LogKeys["NODE_POST_SEND"] = true
	logg.LogKeys["NODE_POST_RECV"] = true
	logg.LogKeys["NODE_STATE"] = true
	ng.SeedRandom()
}

func main() {
	checkerlution.RunTopologyMutatingTrainer()
	// checkerlution.RunPopulationTrainer()
}
