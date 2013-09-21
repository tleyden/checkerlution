
A "checkers bot" which uses [neurgo](https://github.com/tleyden/neurgo) to do it's thinking (or lack thereof).

# How it's modeled

* sensor1: game state array with 32 elements, each of which is:
    * -1.0: opponent king
    * -0.5: opponent piece
    * 0 empty
    * 0.5 our piece
    * 1.0: our king

* sensor2: an available move, which is: 
    * start_location(normalized to be between -1 and 1)
    * is_king(-1: false, 1: true)
    * final_location(-1 and 1)
    * will_be_king(-1: false, 1: true) 
    * amt_would_capture(-1: none, 0: 1 piece, 1: 2 or more pieces)

* actuator: outputs a scalar value representing the confidence in the available move

There is a loop which presents each move to the network, and the move which has the highest confidence wins.

