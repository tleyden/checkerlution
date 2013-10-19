package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	"time"
)

type CheckerlutionScape struct {
	thinker        *Checkerlution
	team           cbot.TeamType
	syncGatewayUrl string
	feedType       cbot.FeedType
}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	logg.LogTo("MAIN", "scape.Fitness() called, playing checkers game")

	fitnessVals := []float64{}

	for i := 0; i < 5; i++ {

		cortex.Init()

		// play a game of checkers
		scape.thinker.StartWithCortex(cortex, scape.team)
		game := cbot.NewGame(scape.team, scape.thinker)
		game.SetServerUrl(scape.syncGatewayUrl)
		game.SetFeedType(scape.feedType)
		game.GameLoop()
		logg.LogTo("MAIN", "gameLoop finished")

		// get result (TODO: this feels clunky, just call thinker.calcFitness() here)
		latestFitness := scape.thinker.latestFitnessScore
		logg.LogTo("MAIN", "checkerlution scape fitness: %v", latestFitness)
		fitnessVals = append(fitnessVals, latestFitness)

		// wait until the game number increments, otherwise on the
		// next callback to this method, we'll jump into a game which
		// is already over.
		logg.LogTo("MAIN", "checkerlution WaitForNextGame()")
		game.WaitForNextGame()

		cortex.Shutdown()

	}

	logg.LogTo("MAIN", "completed game set, calc average fitness")

	total := 0.0
	for _, fitnessVal := range fitnessVals {
		total += fitnessVal
	}
	fitness = total / float64(len(fitnessVals))
	logg.LogTo("MAIN", "checkerlution avg fitness game set: %v", fitness)

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Dumping cortex to %v", filename)
	cortex.MarshalJSONToFile(filename)

	return

}

func (scape *CheckerlutionScape) SetThinker(thinker *Checkerlution) {
	scape.thinker = thinker
}

func (scape *CheckerlutionScape) SetSyncGatewayUrl(syncGatewayUrl string) {
	scape.syncGatewayUrl = syncGatewayUrl
}

func (scape *CheckerlutionScape) SetTeam(team cbot.TeamType) {
	scape.team = team
}

func (scape *CheckerlutionScape) SetFeedType(feedType cbot.FeedType) {
	scape.feedType = feedType
}
