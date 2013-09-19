package checkerlution

import (
	"encoding/json"
	"github.com/couchbaselabs/go.assert"
	"github.com/couchbaselabs/logg"
	"testing"
)

func TestFetchNewGameDocument(t *testing.T) {

	logg.LogKeys["MAIN"] = true
	client := Client{}
	gameState, possibleMoves := client.FetchNewGameDocument()
	logg.LogTo("MAIN", "gameState: %v", gameState)
	logg.LogTo("MAIN", "possibleMoves: %v", possibleMoves)

}

func TestExtractGameRevision(t *testing.T) {
	jsonString := `{"results":[
{"seq":"*:19","id":"user:B2A54597-94CF-43FA-AA96-15DF05322BDF","changes":[{"rev":"4-8a359bbf81184852079617e14b136fa6"}]}
,{"seq":"*:64","id":"vote:EC4861B1-85D9-48D3-8E5D-E76E02478CDC","changes":[{"rev":"1-cef0ead6f6550d97383685d914eab5d7"}]}
,{"seq":"*:1047","id":"votes:checkers","changes":[{"rev":"9-da14169180d383a0782332fc3db7ea3a"}]}
,{"seq":"*:1065","id":"game:checkers","changes":[{"rev":"1033-b8b5a024007b83d45f078c44999b93ed"}]}
],
"last_seq":"*:1065"}`

	var data GenericMap
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		logg.LogPanic("%v", err)
	}

	client := Client{}
	gameRev := client.extractGameRevision(data)
	assert.Equals(t, gameRev, "1033-b8b5a024007b83d45f078c44999b93ed")

}
