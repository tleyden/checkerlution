package checkerlution

type Votes struct {
	Moves  []VoteMove `json:"moves"`
	TeamId int        `json:"team"`
	GameId int        `json:"game"`
	Count  int        `json:"count"`
	Turn   int        `json:"turn"`
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
	/*
		// TODO locations := []int{ move.validMove.}
		voteMove := VoteMove{
			GameId: votes.GameId,
			Count:  votes.Count,
			// TODO!! PieceId: move.PieceId,
			TeamId: votes.TeamId,
			Turn:   votes.Turn,
			// TODO: Locations:
		}
	*/

}
