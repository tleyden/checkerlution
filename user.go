package checkerlution

type User struct {
	Id         string `json:"_id"`
	Rev        string `json:"_rev"`
	TeamId     int    `json:"team"`
	GameNumber int    `json:"game"`
}
