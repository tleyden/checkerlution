package checkerlution

import (
	"fmt"
	"github.com/couchbaselabs/logg"
	cbot "github.com/tleyden/checkers-bot"
	ng "github.com/tleyden/neurgo"
)

type CheckerlutionScape struct {
	team                  cbot.TeamType
	syncGatewayUrl        string
	feedType              cbot.FeedType
	randomDelayBeforeMove int
	fitnessHistory        map[string]float64
}

func (scape *CheckerlutionScape) FitnessAgainst(cortex *ng.Cortex, opponentCortex *ng.Cortex) (fitness float64) {

	logg.LogTo("CHECKERLUTION_SCAPE", "FitnessAgainst cortex: %v vs opponent: %v", cortex.NodeId.UUID, opponentCortex.NodeId.UUID)

	if cortex == opponentCortex {
		logg.LogPanic("Cannot calculate fitnesss between cortex %p and itself %p", cortex, opponentCortex)
	}

	savedFitness, isPresent := scape.lookupFitnessHistory(cortex, opponentCortex)
	if isPresent {
		fitness = savedFitness
		logg.LogTo("CHECKERLUTION_SCAPE", "Fitness from history (us %v): %v", cortex.NodeId.UUID, fitness)
		return
	}

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
	opponentTeam := cbot.RED_TEAM
	if scape.team == cbot.RED_TEAM {
		opponentTeam = cbot.BLUE_TEAM
	}
	thinkerOpponent.StartWithCortex(opponentCortex, opponentTeam)
	gameOpponent := cbot.NewGame(opponentTeam, thinkerOpponent)
	gameOpponent.SetServerUrl(scape.syncGatewayUrl)
	gameOpponent.SetFeedType(scape.feedType)
	gameOpponent.SetDelayBeforeMove(scape.randomDelayBeforeMove)

	// run both game loops and wait for both to finish
	logg.LogTo("CHECKERLUTION_SCAPE", "Started game: %v vs %v", cortex.NodeId.UUID, opponentCortex.NodeId.UUID)
	games := []*cbot.Game{game, gameOpponent}
	scape.runGameLoops(games)
	logg.LogTo("CHECKERLUTION_SCAPE", "Game finished after %v turns", game.Turn())

	fitness = thinker.latestFitnessScore
	fitnessOpponent := thinkerOpponent.latestFitnessScore

	logg.LogTo("CHECKERLUTION_SCAPE", "Fitness (us %v): %v", cortex.NodeId.UUID, fitness)
	logg.LogTo("CHECKERLUTION_SCAPE", "Fitness (opponent %v): %v", opponentCortex.NodeId.UUID, fitnessOpponent)

	scape.recordFitness(cortex, fitness, opponentCortex, fitnessOpponent)

	// wait until the game number increments, otherwise on the
	// next callback to this method, we'll jump into a game which
	// is already over.
	logg.LogTo("CHECKERLUTION_SCAPE", "Wait For Next Game ..")
	game.WaitForNextGame()
	gameOpponent.WaitForNextGame()
	logg.LogTo("CHECKERLUTION_SCAPE", "Done waiting For Next Game ..")

	cortex.Shutdown()
	opponentCortex.Shutdown()

	return
}

func (scape *CheckerlutionScape) Fitness(cortex *ng.Cortex) (fitness float64) {

	logg.LogTo("CHECKERLUTION_SCAPE", "CheckerlutionScape Fitness() called, create random opponent")

	// create an opponent
	thinker := &Checkerlution{}
	thinker.CreateNeurgoCortex()
	thinker.Cortex()

	return scape.FitnessAgainst(cortex, thinker.Cortex())

}

func (scape *CheckerlutionScape) runGameLoops(games []*cbot.Game) {

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

func (scape *CheckerlutionScape) recordFitness(cortex *ng.Cortex, fitness float64, opponentCortex *ng.Cortex, fitnessOpponent float64) {

	if scape.fitnessHistory == nil {
		scape.fitnessHistory = make(map[string]float64)
	}

	// record the fitness score of us vs them
	key := fmt.Sprintf("%v-%v", cortex.NodeId.UUID, opponentCortex.NodeId.UUID)
	scape.fitnessHistory[key] = fitness

}

func (scape *CheckerlutionScape) lookupFitnessHistory(cortex *ng.Cortex, opponentCortex *ng.Cortex) (fitness float64, isPresent bool) {

	key := fmt.Sprintf("%v-%v", cortex.NodeId.UUID, opponentCortex.NodeId.UUID)
	fitness, isPresent = scape.fitnessHistory[key]
	return

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
