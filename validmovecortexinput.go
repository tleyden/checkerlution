package checkerlution

import (
	ng "github.com/tleyden/neurgo"
)

type ValidMoveCortexInput struct {
	validMove       ValidMove
	startLocation   float64
	isCurrentlyKing float64
	endLocation     float64
	willBecomeKing  float64
	captureValue    float64 // -1: 0 capture, 0: 1 capture, 1: 2+ captures
}

func NewValidMoveCortexInput(validMove ValidMove, piece Piece) ValidMoveCortexInput {

	// enhance the validMove from some information from the piece
	validMove.StartLocation = piece.Location
	validMove.PieceId = piece.PieceId

	moveInput := ValidMoveCortexInput{}
	moveInput.validMove = validMove

	// end location
	if len(validMove.Locations) > 0 {
		// the locations array for a valid move can have multiple values,
		// because of "double jumps".  as s simplification, just look at the
		// last jump
		lastIndex := len(validMove.Locations) - 1
		endLocation1Based := validMove.Locations[lastIndex]
		endLocation0Based := endLocation1Based - 1
		endLocation := float64(endLocation0Based)
		moveInput.endLocation = moveInput.normalize(endLocation)
	}

	// is king
	switch piece.King {
	case true:
		moveInput.isCurrentlyKing = 1.0
	case false:
		moveInput.isCurrentlyKing = -1.0
	}

	// start location
	startLocation0Based := piece.Location - 1
	startLocation := float64(startLocation0Based)
	moveInput.startLocation = moveInput.normalize(startLocation)

	// will become king
	switch validMove.King {
	case true:
		moveInput.willBecomeKing = 1.0
	case false:
		moveInput.willBecomeKing = -1.0
	}

	// capture value
	switch len(validMove.Captures) {
	case 0:
		moveInput.captureValue = -1.0
	case 1:
		moveInput.captureValue = 0.0
	default:
		moveInput.captureValue = 1.0
	}

	return moveInput
}

func (move ValidMoveCortexInput) Equals(other ValidMoveCortexInput) bool {

	return move.startLocation == other.startLocation &&
		move.endLocation == other.endLocation &&
		move.captureValue == other.captureValue &&
		move.isCurrentlyKing == other.isCurrentlyKing &&
		move.willBecomeKing == other.willBecomeKing

}

func (move ValidMoveCortexInput) VectorRepresentation() []float64 {
	vector := make([]float64, 5)
	vector[0] = move.startLocation
	vector[1] = move.isCurrentlyKing
	vector[2] = move.endLocation
	vector[3] = move.willBecomeKing
	vector[4] = move.captureValue
	return vector
}

func (move ValidMoveCortexInput) normalize(value float64) float64 {
	normalizeParams := move.getNormalizeParams()
	return ng.NormalizeInRange(normalizeParams, value)
}

func (move ValidMoveCortexInput) getNormalizeParams() ng.NormalizeParams {
	params := ng.NormalizeParams{
		SourceRangeStart: 0.0,
		SourceRangeEnd:   31.0,
		TargetRangeStart: -1.0,
		TargetRangeEnd:   1.0,
	}
	return params
}
