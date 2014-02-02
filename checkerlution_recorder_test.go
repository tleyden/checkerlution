package checkerlution

import (
	"github.com/couchbaselabs/logg"
)

func init() {
	logg.LogKeys["CHECKERLUTION"] = true
	logg.LogKeys["TEST"] = true
}
