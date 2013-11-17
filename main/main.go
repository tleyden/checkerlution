package main

import (
	_ "expvar"
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
	"net/http"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["CHECKERLUTION_SCAPE"] = true
	logg.LogKeys["NEURGO"] = false
	logg.LogKeys["SENSOR_SYNC"] = false
	logg.LogKeys["ACTUATOR_SYNC"] = false
	logg.LogKeys["NODE_PRE_SEND"] = false
	logg.LogKeys["NODE_POST_SEND"] = false
	logg.LogKeys["NODE_POST_RECV"] = false
	logg.LogKeys["NODE_STATE"] = false
	ng.SeedRandom()
}

func main() {

	// run a webserver in order to view expvar output
	// at http://localhost:8080/debug/vars
	go http.ListenAndServe(":8080", nil)

	trainer := &checkerlution.CheckerlutionTrainer{}

	checkerlution.RegisterHandlers(trainer)

	// checkerlution.RunTopologyMutatingTrainer()
	trainer.RunPopulationTrainer()

}
