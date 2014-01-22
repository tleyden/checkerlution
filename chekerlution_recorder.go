package checkerlution

import (
	"github.com/couchbaselabs/logg"
	"github.com/tleyden/go-couch"
	ng "github.com/tleyden/neurgo"
	nv "github.com/tleyden/neurvolve"
)

type CheckerlutionRecorder struct {
	syncGatewayUrl string
	db             couch.Database
	population     Population
}

// This assumes the population does not already exist in the database
// and so it will create it
func NewRecorder(syncGatewayUrl string, population Population) *CheckerlutionRecorder {
	recorder := &CheckerlutionRecorder{
		syncGatewayUrl: syncGatewayUrl,
		population:     population,
	}
	recorder.InitDbConnection()

	newId, rev, err := recorder.db.InsertWith(population, population.name)
	if err != nil {
		logg.LogTo("CHECKERLUTION", "Error saving population document: %v", err.Error())
	} else {
		logg.LogTo("CHECKERLUTION", "Saved empty population document.  Id: %v, Rev: %v, Err: %v", newId, rev, err)
		recorder.population.rev = rev
	}

	return recorder
}

func (r *CheckerlutionRecorder) InitDbConnection() {
	db, error := couch.Connect(r.syncGatewayUrl)
	if error != nil {
		logg.LogPanic("Error connecting to %v: %v", r.syncGatewayUrl, error)
	}
	r.db = db
}

func (r *CheckerlutionRecorder) AddGeneration(evaldCortexes []nv.EvaluatedCortex) {

	logg.LogTo("CHECKERLUTION", "AddGeneration called.  Rev: %v", r.population.rev)
	agents := []Agent{}
	for _, evaldCortex := range evaldCortexes {
		parent_id := evaldCortex.Cortex.NodeId.UUID
		agent := NewAgent(evaldCortex.Cortex, parent_id)
		agents = append(agents, *agent)
	}
	generationNumber := r.population.NextGenerationNumber()
	generation := NewGeneration(generationNumber, agents)
	(*generation).state = "initial"
	r.population.AddGeneration(*generation)

	r.Save()

}

func (r *CheckerlutionRecorder) AddFitnessScore(score float64, cortex *ng.Cortex, opponent *ng.Cortex) {

	// find the "current" generation where this game should be added
	generation := r.population.CurrentGeneration()

	// figure out which was winner based on score
	winner_id := cortex.NodeId.UUID
	if score == 0.0 {
		winner_id = "draw"
	} else if score < 0 {
		winner_id = opponent.NodeId.UUID
	}

	// create new game
	game := Game{
		red_player_id:  cortex.NodeId.UUID,
		blue_player_id: opponent.NodeId.UUID,
		winner_id:      winner_id,
	}

	// add to slice of games for that generation
	generation.AddGame(game)

	r.Save()

}

func (r *CheckerlutionRecorder) Save() {
	rev, err := r.db.EditWith(
		r.population,
		r.population.name,
		r.population.rev,
	)
	if err != nil {
		logg.LogPanic("Error adding generation: %v", err.Error())
	}
	r.population.rev = rev
	logg.LogTo("CHECKERLUTION", "Saved population document.  Rev: %v", rev)

}
