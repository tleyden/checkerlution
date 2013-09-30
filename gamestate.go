package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
)

// data structure that corresponds to the checkers:game json doc
type GameState struct {
	Teams        []Team `json:"teams"`
	ActiveTeam   int    `json:"activeTeam"`
	Number       int    `json:"number"`
	Turn         int    `json:"turn"`
	MoveInterval int    `json:"moveInterval"`
}

func NewGameStateFromString(jsonString string) GameState {
	gameState := &GameState{}
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, gameState)
	if err != nil {
		logg.LogError(err)
	}
	return *gameState
}

type Piece struct {
	Location   int         `json:"location"`
	King       bool        `json:"king"`
	Captured   bool        `json:"captured"`
	ValidMoves []ValidMove `json:"validMoves"`
	PieceId    int
}

type Team struct {
	Score            int     `json:"score"`
	ParticipantCount int     `json:"participantCount"`
	Pieces           []Piece `json:"pieces"`
}

type ValidMove struct {
	Locations     []int     `json:"locations"`
	Captures      []Capture `json:"captures"`
	King          bool      `json:"king"`
	PieceId       int
	StartLocation int
}

type Capture struct {
	TeamID  int `json:"team"`
	PieceId int `json:"piece"`
}

func (t Team) AllValidMoves() (validMoves []ValidMove) {
	validMoves = make([]ValidMove, 0)

	for pieceIndex, piece := range t.Pieces {
		for _, validMove := range piece.ValidMoves {
			// enhance the validMove from some information
			// from the piece
			validMove.StartLocation = piece.Location
			validMove.PieceId = pieceIndex

			validMoves = append(validMoves, validMove)
		}
	}
	return
}
