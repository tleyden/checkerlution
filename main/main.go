package main

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/checkerlution"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
	"net/http"
	"net/url"
)

func init() {
	logg.LogKeys["MAIN"] = true
	logg.LogKeys["DEBUG"] = true
	ng.SeedRandom()
}

var shouldUseProxy bool = false

func main() {

	if shouldUseProxy {
		useProxy()
	}

	// create a checkerlution instance just to create a cortex (kludgy)
	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.TRAINING_MODE)
	thinker.CreateNeurgoCortex()
	cortex := thinker.Cortex()

	// setup the scape
	scape := &checkerlution.CheckerlutionScape{}
	scape.SetThinker(thinker)

	// create a stochastic hill climber
	shc := &nv.StochasticHillClimber{
		FitnessThreshold:           50,
		MaxIterationsBeforeRestart: 5,
		MaxAttempts:                1,
	}
	cortexTrained, succeeded := shc.TrainScape(cortex, scape)
	if succeeded {
		logg.LogTo("MAIN", "Training succeeded")
	} else {
		logg.LogTo("MAIN", "Training Failed")
	}

	cortexTrained.MarshalJSONToFile("/tmp/checkerlution.json")

}

func useProxy() {
	proxyUrl, err := url.Parse("http://localhost:8888")
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	if err != nil {
		panic("proxy issue")
	}

}
