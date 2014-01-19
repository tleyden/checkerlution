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

	newId, newRevision, err := recorder.db.InsertWith(population, population.name)
	if err != nil {
		logg.LogTo("CHECKERLUTION", "Error saving population document: %v", err.Error())
	} else {
		logg.LogTo("CHECKERLUTION", "Saved empty population document.  Id: %v, Rev: %v, Err: %v", newId, newRevision, err)
		recorder.population.revision = newRevision
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

func (r CheckerlutionRecorder) AddGeneration(evaldCortexes []nv.EvaluatedCortex) {

	agents := []Agent{}
	for _, evaldCortex := range evaldCortexes {
		agent := NewAgent(evaldCortex.Cortex)
		agents = append(agents, *agent)
	}
	generationNumber := r.population.NextGenerationNumber()
	generation := NewGeneration(generationNumber, agents)
	r.population.AddGeneration(*generation)
	r.population.test = "foo"
	newRevision, err := r.db.EditWith(r.population, r.population.name, r.population.revision)
	if err != nil {
		logg.LogPanic("Error adding generation: %v", err.Error())
	}
	logg.LogTo("CHECKERLUTION", "Saved population document.  Rev: %v", newRevision)

}

func (r CheckerlutionRecorder) AddFitnessScore(score float64, cortex *ng.Cortex, opponent *ng.Cortex) {

}
