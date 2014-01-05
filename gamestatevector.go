package checkerlution

import (
	cbot "github.com/tleyden/checkers-bot"
	core "github.com/tleyden/checkers-core"
)

const OPPONENT_KING = -1.0
const OPPONENT_PIECE = -0.7
const EMPTY_SQUARE = 0.0
const OUR_PIECE = 0.7
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

func (v *GameStateVector) loadFromBoard(board core.Board, ourPlayer core.Player) {
	v2 := *v
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			loc := core.NewLocation(row, col)
			vectorIndex1Based := cbot.ExportCoreLocation(loc)
			if vectorIndex1Based == -1 {
				// invalid square - ignore
				continue
			}
			vectorIndex := vectorIndex1Based - 1 // convert to 0-based
			piece := board.PieceAt(loc)
			switch ourPlayer {
			case core.BLACK_PLAYER:
				switch piece {
				case core.BLACK:
					v2[vectorIndex] = OUR_PIECE
				case core.BLACK_KING:
					v2[vectorIndex] = OUR_KING
				case core.RED:
					v2[vectorIndex] = OPPONENT_PIECE
				case core.RED_KING:
					v2[vectorIndex] = OPPONENT_KING
				default:
					v2[vectorIndex] = EMPTY_SQUARE
				}

			default:
				switch piece {
				case core.BLACK:
					v2[vectorIndex] = OPPONENT_PIECE
				case core.BLACK_KING:
					v2[vectorIndex] = OPPONENT_KING
				case core.RED:
					v2[vectorIndex] = OUR_PIECE
				case core.RED_KING:
					v2[vectorIndex] = OUR_KING
				default:
					v2[vectorIndex] = EMPTY_SQUARE
				}

			}

		}
	}

}

// Deprecated - remove
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
