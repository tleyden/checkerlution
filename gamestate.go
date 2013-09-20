package checkerlution

// data structure that corresponds to the checkers:game json doc
type GameState struct {
	Teams      []Team `json:"teams"`
	ActiveTeam int    `json:"activeTeam"`
	Number     int    `json:"number"`
}

type Piece struct {
	Location   int         `json:"location"`
	King       bool        `json:"king"`
	Captured   bool        `json:"captured"`
	ValidMoves []ValidMove `json:"validMoves"`
}

type Team struct {
	Score            int     `json:"score"`
	ParticipantCount int     `json:"participantCount"`
	Pieces           []Piece `json:"pieces"`
}

type ValidMove struct {
	Locations []int     `json:"locations"`
	Captures  []Capture `json:"captures"`
	King      bool      `json:"king"`
}

type Capture struct {
	TeamID  int `json:"team"`
	PieceId int `json:"piece"`
}
