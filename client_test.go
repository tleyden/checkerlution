package checkerlution

import (
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
