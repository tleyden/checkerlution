package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
	"time"
)

type CheckerlutionScape struct {
	thinker               *Checkerlution
	team                  cbot.TeamType
	syncGatewayUrl        string
	feedType              cbot.FeedType
	randomDelayBeforeMove int
}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex, opponentCortex *ng.Cortex) (fitness float64) {

	cortex.Init()
	opponentCortex.Init()

	// setup checkers game for cortex
	thinker := &checkerlution.Checkerlution{}
	thinker.SetMode(checkerlution.TRAINING_MODE)
	thinker.StartWithCortex(cortex, scape.team)
	game := cbot.NewGame(scape.team, thinker)
	game.SetServerUrl(scape.syncGatewayUrl)
	game.SetFeedType(scape.feedType)
	game.SetDelayBeforeMove(scape.randomDelayBeforeMove)

	// setup checkers game for opponent
	thinkerOpponent := &checkerlution.Checkerlution{}
	thinkerOpponent.SetMode(checkerlution.TRAINING_MODE)
	oppontentTeam := cbot.RED_TEAM
	if scape.team == cbot.RED_TEAM {
		opponentTeam = cbot.BLUE_TEAM
	}
	thinkerOpponent.StartWithCortex(opponentCortex, oppontentTeam)
	gameOpponent := cbot.NewGame(oppontentTeam, thinkerOpponent)
	gameOpponent.SetServerUrl(scape.syncGatewayUrl)
	gameOpponent.SetFeedType(scape.feedType)
	gameOpponent.SetDelayBeforeMove(scape.randomDelayBeforeMove)

	// run both game loops and wait for both to finish
	games := []*cbot.Game{game, gameOpponent}
	scape.runGameLoops(games)

	fitness = thinker.latestFitnessScore
	logg.LogTo("MAIN", "Fitness: %v", fitness)

	// wait until the game number increments, otherwise on the
	// next callback to this method, we'll jump into a game which
	// is already over.
	logg.LogTo("MAIN", "Wait For Next Game ..")
	game.WaitForNextGame()

	logg.LogTo("MAIN", "Wait For Next Opponent Game ..")
	gameOpponent.WaitForNextGame()

	cortex.Shutdown()
	opponentCortex.Shutdown()

	return

}

func (scape *CheckerlutionScape) runGameLoops(games []*cbot.Game) {

	resultChannel := make(chan bool)
	runGameLoop := func(game *cbot.Game, result chan bool) {
		game.GameLoop()
		result <- true
	}
	for i, game := range games {
		go runGameLoop(game, resultChannel)
	}
	for i, game := range games {
		<-resultChannel
	}

}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	logg.LogTo("DEBUG", "scape.Fitness() called, playing checkers game")

	fitnessVals := []float64{}

	for i := 0; i < 7; i++ {

		cortex.Init()

		// play a game of checkers
		scape.thinker.StartWithCortex(cortex, scape.team)
		game := cbot.NewGame(scape.team, scape.thinker)
		game.SetServerUrl(scape.syncGatewayUrl)
		game.SetFeedType(scape.feedType)
		game.SetDelayBeforeMove(scape.randomDelayBeforeMove)
		game.GameLoop()
		logg.LogTo("DEBUG", "gameLoop finished")

		// get result (TODO: this feels clunky, just call thinker.calcFitness() here)
		latestFitness := scape.thinker.latestFitnessScore
		logg.LogTo("MAIN", "Fitness: %v", latestFitness)
		fitnessVals = append(fitnessVals, latestFitness)

		// wait until the game number increments, otherwise on the
		// next callback to this method, we'll jump into a game which
		// is already over.
		logg.LogTo("MAIN", "Wait For Next Game ..")
		game.WaitForNextGame()

		cortex.Shutdown()

	}

	logg.LogTo("DEBUG", "Game Set Finished")

	total := 0.0
	for _, fitnessVal := range fitnessVals {
		total += fitnessVal
	}
	fitness = total / float64(len(fitnessVals))
	logg.LogTo("MAIN", "Average Fitness for Game Set: %v", fitness)

	filename := fmt.Sprintf("/tmp/checkerlution-%v.json", time.Now().Unix())
	logg.LogTo("MAIN", "Saving Cortex snapshot to %v", filename)
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

func (scape *CheckerlutionScape) SetRandomDelayBeforeMove(delay int) {
	scape.randomDelayBeforeMove = delay
}
