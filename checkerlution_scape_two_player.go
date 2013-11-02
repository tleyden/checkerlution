package checkerlution

import (
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

type CheckerlutionScapeTwoPlayer struct {
	team                  cbot.TeamType
	syncGatewayUrl        string
	feedType              cbot.FeedType
	randomDelayBeforeMove int
}

func (scape *CheckerlutionScapeTwoPlayer) Fitness(cortex *ng.Cortex, opponentCortex *ng.Cortex) (fitness float64) {

	cortex.Init()
	opponentCortex.Init()

	// setup checkers game for cortex
	thinker := &Checkerlution{}
	thinker.SetMode(TRAINING_MODE)
	thinker.StartWithCortex(cortex, scape.team)
	game := cbot.NewGame(scape.team, thinker)
	game.SetServerUrl(scape.syncGatewayUrl)
	game.SetFeedType(scape.feedType)
	game.SetDelayBeforeMove(scape.randomDelayBeforeMove)

	// setup checkers game for opponent
	thinkerOpponent := &Checkerlution{}
	thinkerOpponent.SetMode(TRAINING_MODE)
	var opponentTeam cbot.TeamType // TODO: why can't just use := syntax here?
	opponentTeam = cbot.RED_TEAM
	if scape.team == cbot.RED_TEAM {
		opponentTeam = cbot.BLUE_TEAM
	}
	thinkerOpponent.StartWithCortex(opponentCortex, opponentTeam)
	gameOpponent := cbot.NewGame(opponentTeam, thinkerOpponent)
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

func (scape *CheckerlutionScapeTwoPlayer) runGameLoops(games []*cbot.Game) {

	resultChannel := make(chan bool)
	runGameLoop := func(game *cbot.Game, result chan bool) {
		game.GameLoop()
		result <- true
	}
	for _, game := range games {
		go runGameLoop(game, resultChannel)
	}
	for _, _ = range games {
		<-resultChannel
	}

}

func (scape *CheckerlutionScapeTwoPlayer) SetSyncGatewayUrl(syncGatewayUrl string) {
	scape.syncGatewayUrl = syncGatewayUrl
}

func (scape *CheckerlutionScapeTwoPlayer) SetTeam(team cbot.TeamType) {
	scape.team = team
}

func (scape *CheckerlutionScapeTwoPlayer) SetFeedType(feedType cbot.FeedType) {
	scape.feedType = feedType
}

func (scape *CheckerlutionScapeTwoPlayer) SetRandomDelayBeforeMove(delay int) {
	scape.randomDelayBeforeMove = delay
}
