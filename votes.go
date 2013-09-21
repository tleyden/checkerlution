package checkerlution

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
