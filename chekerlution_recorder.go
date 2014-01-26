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
func NewRecorder(syncGatewayUrl string, populationName string) *CheckerlutionRecorder {
	recorder := &CheckerlutionRecorder{
		syncGatewayUrl: syncGatewayUrl,
	}
	recorder.InitDbConnection()

	recorder.createOrRetrievePopulation(populationName)

	return recorder
}

func (r *CheckerlutionRecorder) createOrRetrievePopulation(populationName string) {

	population := &Population{}

	// try to fetch existing population
	logg.LogTo("CHECKERLUTION", "Looking up: %v", populationName)
	error := r.db.Retrieve(populationName, population)

	if error != nil {
		// if does not exist, create new one
		logg.LogTo("CHECKERLUTION", "Could not find: %v", populationName, error)

		population.name = populationName
		newId, rev, err := r.db.InsertWith(population, populationName)
		if err != nil {
			logg.LogPanic("Error saving population document: %v", err.Error())
		} else {
			logg.LogTo("CHECKERLUTION", "Saved empty population document.  Id: %v, Rev: %v, Err: %v", newId, rev, err)
			population.rev = rev

		}

	} else {
		logg.LogTo("CHECKERLUTION", "Found: %v", *population)
		logg.LogTo("CHECKERLUTION", "Rev: %v", (*population).rev)

	}

	r.population = *population

}

func (r CheckerlutionRecorder) GetLatestGenerationCortexes() []*ng.Cortex {
	latestGeneration := Generation{}
	for _, generation := range r.population.generations {
		latestGeneration = generation
	}

	cortexes := []*ng.Cortex{}
	for _, agent := range latestGeneration.agents {
		cortex := LoadCortex(agent.cortex_id, r.db)
		cortexes = append(cortexes, cortex)
	}
	return cortexes

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
		parent_id := evaldCortex.ParentId
		if len(parent_id) == 0 {
			parent_id = evaldCortex.Cortex.NodeId.UUID
		}
		r.SaveCortex(evaldCortex.Cortex)
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

func (r *CheckerlutionRecorder) SaveCortex(cortex *ng.Cortex) {

	newId, rev, err := r.db.InsertWith(cortex, cortex.NodeId.UUID)
	if err != nil {
		logg.LogTo("CHECKERLUTION", "Error saving cortex: %v, already exists?", err.Error())
	} else {
		logg.LogTo("CHECKERLUTION", "Saved cortex.  Id: %v, Rev: %v, Err: %v", newId, rev, err)
	}

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

func LoadCortex(cortexId string, db couch.Database) *ng.Cortex {
	cortex := &ng.Cortex{}
	error := db.Retrieve(cortexId, cortex)
	if error != nil {
		logg.LogPanic("Could not find cortex: %v", cortexId, error)
	}
	cortex.LinkNodesToCortex()
	return cortex
}
