package main

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = false
	logg.LogKeys["NEURGO"] = false
	logg.LogKeys["NODE_PRE_SEND"] = false
	logg.LogKeys["NODE_POST_SEND"] = false
	logg.LogKeys["NODE_POST_RECV"] = false
	logg.LogKeys["NODE_STATE"] = false
	ng.SeedRandom()
}

func main() {
	runTopologyMutatingTrainer()
}
