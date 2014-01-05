package checkerlution

import (
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

func init() {
	logg.LogKeys["DEBUG"] = true
	logg.LogKeys["CHECKERLUTION"] = true
	logg.LogKeys["CHECKERLUTION_SCAPE"] = true
	ng.SeedRandom()
}
