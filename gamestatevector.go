package checkerlution

import (
	cbot "github.com/tleyden/checkers-bot"
)

const OPPONENT_KING = -1.0
const OPPONENT_PIECE = -0.5
const EMPTY_SQUARE = 0.0
const OUR_PIECE = 0.5
const OUR_KING = 1.0

// game state array with 32 elements, possible values:
// OPPONENT_KING
// OPPONENT_PIECE
// EMPTY_SQUARE
// OUR_PIECE
// OUR_KING
type GameStateVector []float64

func NewGameStateVector() GameStateVector {

	gameStateVector := make(GameStateVector, 32)
	for i := 0; i < 32; i++ {
		gameStateVector[i] = EMPTY_SQUARE
	}

	return gameStateVector

}

func (v *GameStateVector) loadFromGameState(gameState cbot.GameState, ourTeamId cbot.TeamType) {

	v2 := *v
	for teamID, team := range gameState.Teams {
		isOurTeam := (teamID == int(ourTeamId))
		for _, piece := range team.Pieces {
			vectorIndex := piece.Location - 1

			switch isOurTeam {
			case true:
				switch piece.King {
				case true:
					v2[vectorIndex] = OUR_KING
				case false:
					v2[vectorIndex] = OUR_PIECE
				}
			case false:
				switch piece.King {
				case true:
					v2[vectorIndex] = OPPONENT_KING
				case false:
					v2[vectorIndex] = OPPONENT_PIECE
				}
			}
		}
	}

}
