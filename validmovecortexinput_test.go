package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
	"testing"
)

func init() {
	logg.LogNoColor()
	logg.LogKeys["TEST"] = true
}

func TestNewValidMoveCortexInput(t *testing.T) {

	jsonString := `{"locations":[1],"captures":[{"team":1,"piece":11}],"king":true}`

	validMovePtr := &ValidMove{}
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, validMovePtr)
	if err != nil {
		logg.LogError(err)
	}

	validMove := *validMovePtr
	vmCortexInput := NewValidMoveCortexInput(validMove)
	endLocationNormalized := vmCortexInput.endLocation
	expected := -1.0
	logg.LogTo("TEST", "endLocationNormalized: %v", endLocationNormalized)

	assert.True(t, ng.EqualsWithMaxDelta(endLocationNormalized, expected, 0.1))

	logg.LogTo("TEST", "vmCortexInput: %v", vmCortexInput)

}
