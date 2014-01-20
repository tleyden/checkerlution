package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	nv "github.com/tleyden/neurvolve"
	"testing"
)

func init() {
	logg.LogKeys["CHECKERLUTION"] = true
}

func TestAddGeneration(t *testing.T) {

	population := Population{name: "foo18"}
	recorder := NewRecorder("http://localhost:4984/checkers", population)

	cortex := nv.BasicCortex()
	evaldCortex := nv.EvaluatedCortex{
		Cortex:  cortex,
		Fitness: 0.0,
	}
	recorder.AddGeneration([]nv.EvaluatedCortex{evaldCortex})

	cortex2 := nv.BasicCortex()
	evaldCortex2 := nv.EvaluatedCortex{
		Cortex:  cortex2,
		Fitness: 0.0,
	}
	recorder.AddGeneration([]nv.EvaluatedCortex{evaldCortex, evaldCortex2})

	assert.True(t, true)
}
