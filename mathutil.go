package checkerlution

import (
	"github.com/couchbaselabs/logg"
	"math/rand"
)

func randomIntInRange(min, max int) int {
	if min == max {
		logg.Log("warn: min==max (%v == %v)", min, max)
		return min
	}
	return rand.Intn(max-min) + min
}

func randomInRange(min, max float64) float64 {

	return rand.Float64()*(max-min) + min
}
