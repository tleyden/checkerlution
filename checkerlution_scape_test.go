package checkerlution

import (
	"github.com/couchbaselabs/go.assert"
	cbot "github.com/tleyden/checkers-bot"
	"testing"
)

func TestLookupFitnessHistory(t *testing.T) {

	checkerlution := &Checkerlution{}
	checkerlution.ourTeamId = cbot.RED_TEAM
	checkerlution.CreateNeurgoCortex()
	cortex := checkerlution.cortex

	checkerlutionOpponent := &Checkerlution{}
	checkerlutionOpponent.ourTeamId = cbot.BLUE_TEAM

	checkerlutionOpponent.CreateNeurgoCortex()
	opponentCortex := checkerlutionOpponent.cortex

	scape := &CheckerlutionScape{}

	fitness := 10.0
	fitnessOpponent := -10.0

	scape.recordFitness(cortex, fitness, opponentCortex, fitnessOpponent)
	retrievedFitness, isPresent := scape.lookupFitnessHistory(cortex, opponentCortex)
	assert.True(t, isPresent)
	assert.Equals(t, retrievedFitness, 10.0)

}
