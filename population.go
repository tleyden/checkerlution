package checkerlution

import (
	"encoding/json"
	ng "github.com/tleyden/neurgo"
)

type Population struct {
	name        string
	rev         string
	generations []Generation
}

func (population Population) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Name        string
			Generations []Generation
		}{
			Name:        population.name,
			Generations: population.generations,
		})
}

func (population Population) NextGenerationNumber() int {
	return len(population.generations)
}

type Agent struct {
	cortex *ng.Cortex
}

func (agent Agent) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Cortex *ng.Cortex
		}{
			Cortex: agent.cortex,
		})
}

func NewAgent(cortex *ng.Cortex) *Agent {
	return &Agent{cortex: cortex}
}

type Generation struct {
	start_time string
	number     int
	state      string
	agents     []Agent
}

func (generation Generation) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Agents []Agent
			Number int
		}{
			Agents: generation.agents,
			Number: generation.number,
		})
}

func NewGeneration(number int, agents []Agent) *Generation {
	return &Generation{
		number: number,
		agents: agents,
	}
}

func (population *Population) AddGeneration(generation Generation) {
	population.generations = append(population.generations, generation)
}

// 	r.population.AddGeneration(generation)

/*
		agent := NewAgent(evaldCortex.Cortex)
		agents := append(agents, agent)
	}
	generationNumber := r.population.NextGenerationNumber()
	generation := NewGeneration(generationNumber, agents)

*/

/*
type Population struct {
	name        string
	generations []Generation
}

func NewPopulation(name string, generation Generation) *Population {
	return &Population{
		name:        name,
		generations: []Generation{generation},
	}
}

type Generation struct {
	start_time string
	number     int
	state      string
	agents     []Agent
}

func NewGeneration(generation []*ng.Cortex) *Generation {
	generation := &Generation{
		number: 0, // TODO: fix this

	}
}

type Agent struct {
	cortex *ng.Cortex
}

func NewAgent(cortex *ng.Cortex) *Agent {

}
*/
