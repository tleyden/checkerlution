package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/logg"
	"io/ioutil"
	"net/http"
)

const CHANGES_FEED_URL = "http://localhost:4984/checkers/_changes?feed=longpoll&timeout=20000"

type GenericMap map[string]interface{}

type Client struct {
}

func (client Client) FetchChangesFeed() (data GenericMap) {

	url := CHANGES_FEED_URL
	resp, fetch_err := http.Get(url)
	logg.LogTo("MAIN", "resp: %v", resp)
	if fetch_err != nil {
		logg.LogPanic("Failed to fetch url: %v.  Err: %v", url, fetch_err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logg.LogPanic("Failed to fetch content from: %v.  Err: %v", url, err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logg.LogPanic("%v", err)
	}
	logg.LogTo("MAIN", "data: %v", data)

	return
}

func (client Client) FetchNewGameDocument() (gameState []float64, possibleMoves []Move) {

	changesFeedMap := client.FetchChangesFeed()
	results := changesFeedMap["results"]
	logg.LogTo("MAIN", "results: %v", results)

	// TODO: this should be
	// - pulled from server
	// - parsed into json
	// - data structs should be extracted from json

	gameState = make([]float64, 32)

	possibleMove1 := Move{
		startLocation:   0,
		isCurrentlyKing: -1,
		endLocation:     1.0,
		willBecomeKing:  -0.5,
		captureValue:    1,
	}

	possibleMove2 := Move{
		startLocation:   1,
		isCurrentlyKing: -0.5,
		endLocation:     0.0,
		willBecomeKing:  0.5,
		captureValue:    0,
	}

	possibleMoves = []Move{possibleMove1, possibleMove2}
	return
}
