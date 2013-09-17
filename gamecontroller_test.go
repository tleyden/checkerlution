package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	"testing"
)

func TestCreateNeurgoCortex(t *testing.T) {
	game := &Game{}
	game.CreateNeurgoCortex()
	cortex := game.cortex
	assert.True(t, cortex != nil)
	assert.True(t, cortex.Sensors != nil)
}
