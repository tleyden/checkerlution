package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	ng "github.com/tleyden/neurgo"
)

type Population struct {
	name        string
	rev         string
	generations []Generation
}

// hack to be able to go back and forth from json
type PopulationJson struct {
	Name        string
	Rev         string `json:"_rev"`
	Generations []Generation
}

func (population Population) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&PopulationJson{
			Name:        population.name,
			Generations: population.generations,
		})
}

func (population *Population) UnmarshalJSON(data []byte) error {

	logg.LogTo("CHECKERLUTION", "")
	populationJson := &PopulationJson{}
	error := json.Unmarshal(data, populationJson)
	if error == nil {
		population.name = populationJson.Name
		population.rev = populationJson.Rev
		population.generations = populationJson.Generations
	}
	return error

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
	cortex_id string
	parent_id string
}

type AgentJson struct {
	Cortex_id string
	Parent_id string
}

func (agent Agent) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&AgentJson{
			Cortex_id: agent.cortex_id,
			Parent_id: agent.parent_id,
		})

}

func (agent *Agent) UnmarshalJSON(data []byte) error {

	agentJson := &AgentJson{}
	error := json.Unmarshal(data, agentJson)
	if error == nil {
		agent.cortex_id = agentJson.Cortex_id
		agent.parent_id = agentJson.Parent_id

	}
	return error

}

func NewAgent(cortex *ng.Cortex, parent_id string) *Agent {
	return &Agent{
		cortex_id: cortex.NodeId.UUID,
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

type GenerationJson struct {
	Start_time string
	Number     int
	State      string
	Agents     []Agent
	Games      []Game
}

func (generation *Generation) AddGame(game Game) {
	generation.games = append(generation.games, game)
}

func (generation *Generation) UnmarshalJSON(data []byte) error {

	generationJson := &GenerationJson{}
	error := json.Unmarshal(data, generationJson)
	if error == nil {
		generation.start_time = generationJson.Start_time
		generation.number = generationJson.Number
		generation.state = generationJson.State
		generation.agents = generationJson.Agents
		generation.games = generationJson.Games
	}
	return error

}

func (generation Generation) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&GenerationJson{
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

type GameJson struct {
	Red_player_id  string
	Blue_player_id string
	Winner_id      string
}

func (game Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&GameJson{
			Red_player_id:  game.red_player_id,
			Blue_player_id: game.blue_player_id,
			Winner_id:      game.winner_id,
		})
}

func (game *Game) UnmarshalJSON(data []byte) error {

	gameJson := &GameJson{}
	error := json.Unmarshal(data, gameJson)
	if error == nil {
		game.red_player_id = gameJson.Red_player_id
		game.blue_player_id = gameJson.Blue_player_id
		game.winner_id = gameJson.Winner_id
	}
	return error

}
