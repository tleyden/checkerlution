package checkerlution

import ()

type Move struct {
	startLocation   float64
	isCurrentlyKing float64
	endLocation     float64
	willBecomeKing  float64
	captureValue    float64 // -1: 0 capture, 0: 1 capture, 1: 2+ captures
}

func (move Move) VectorRepresentation() []float64 {
	vector := make([]float64, 5)
	vector[0] = move.startLocation
	vector[1] = move.isCurrentlyKing
	vector[2] = move.endLocation
	vector[3] = move.willBecomeKing
	vector[4] = move.captureValue
	return vector
}
