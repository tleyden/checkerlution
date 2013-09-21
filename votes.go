package checkerlution

type Votes struct {
	// Id        string                 `json:"_id"`
	// Rev       string                 `json:"_rev"`
	// Revisions map[string]interface{} `json:"_revisions"`
	Channels []interface{} `json:"channels"`
	Moves    []VoteMove    `json:"moves"`
	TeamId   int           `json:"team"`
	GameId   int           `json:"game"`
	Count    int           `json:"count"`
	Turn     int           `json:"turn"`
}

type OutgoingVotes struct {
	Id        string                 `json:"_id"`
	Rev       string                 `json:"_rev"`
	Revisions map[string]interface{} `json:"_revisions"`
	Channels  []interface{}          `json:"channels"`
	Moves     []VoteMove             `json:"moves"`
	TeamId    int                    `json:"team"`
	GameId    int                    `json:"game"`
	Count     int                    `json:"count"`
	Turn      int                    `json:"turn"`
	PieceId   int                    `json:"piece"`
	Locations []int                  `json:"locations"`
}

type VoteMove struct {
	GameId    int   `json:"game"`
	Count     int   `json:"count"`
	Locations []int `json:"locations"`
	PieceId   int   `json:"piece"`
	TeamId    int   `json:"team"`
	Turn      int   `json:"turn"`
}

func (votes *Votes) SetMove(move ValidMoveCortexInput) {

	// TODO: this is actually a bug, because if there is a
	// double jump it will only send the first jump move
	endLocation := move.validMove.Locations[0]

	locations := []int{move.validMove.StartLocation, endLocation}
	voteMove := VoteMove{
		GameId:    votes.GameId,
		Count:     votes.Count,
		PieceId:   move.validMove.PieceId,
		TeamId:    votes.TeamId,
		Turn:      votes.Turn,
		Locations: locations,
	}
	moves := []VoteMove{voteMove}
	votes.Moves = moves

}
