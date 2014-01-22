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

func (population Population) CurrentGenerationNumber() int {
	return len(population.generations) - 1
}

func (population Population) CurrentGeneration() *Generation {
	genNumber := population.CurrentGenerationNumber()
	return &(population.generations[genNumber])
}

type Agent struct {
	cortex    *ng.Cortex
	parent_id string
}

func (agent Agent) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Cortex    *ng.Cortex
			Parent_id string
		}{
			Cortex:    agent.cortex,
			Parent_id: agent.parent_id,
		})
}

func NewAgent(cortex *ng.Cortex, parent_id string) *Agent {
	return &Agent{
		cortex:    cortex,
		parent_id: parent_id,
	}
}

type Generation struct {
	start_time string
	number     int
	state      string
	agents     []Agent
	games      []Game
}

func (generation *Generation) AddGame(game Game) {
	generation.games = append(generation.games, game)
}

func (generation Generation) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Number int
			Agents []Agent
			Games  []Game
		}{
			Number: generation.number,
			Agents: generation.agents,
			Games:  generation.games,
		})
}

func NewGeneration(number int, agents []Agent) *Generation {
	games := []Game{}
	return &Generation{
		number: number,
		agents: agents,
		games:  games,
	}
}

func (population *Population) AddGeneration(generation Generation) {
	population.generations = append(population.generations, generation)
}

type Game struct {
	red_player_id  string
	blue_player_id string
	winner_id      string
}

func (game Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Red_player_id  string
			Blue_player_id string
			Winner_id      string
		}{
			Red_player_id:  game.red_player_id,
			Blue_player_id: game.blue_player_id,
			Winner_id:      game.winner_id,
		})
}