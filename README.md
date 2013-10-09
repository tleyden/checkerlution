[![Build Status](https://drone.io/github.com/tleyden/checkerlution/status.png)](https://drone.io/github.com/tleyden/checkerlution/latest)

A checkers bot implementation for [checkers-bot](https://github.com/tleyden/checkers-bot) which uses [neurgo](https://github.com/tleyden/neurgo) to do it's thinking (or lack thereof).

# Architecture

![architecture png](http://cl.ly/image/3q243W3w3900/Screen%20Shot%202013-10-08%20at%2010.35.57%20PM.png)

# Install pre-requisites

* Go 1.1 or later

# Install checkerlution

```
$ go get github.com/tleyden/checkerlution
$ go get github.com/couchbaselabs/go.assert
```
# Validate installation - run tests

```
$ cd $GOPATH/src/github.com/tleyden/checkerlution
$ go test -v
```

# Install Pre-requisites

* [Couchbase Server](http://www.couchbase.com/download)

* [Sync Gateway](https://github.com/couchbase/sync_gateway)

* [Checkers Overlord](https://github.com/apage43/checkers-overlord)

* [Checkers-iOS](https://github.com/couchbaselabs/Checkers-iOS)

Checkers-iOS is not strictly required, but very useful in order to view the game.

The other way to install [Checkers-iOS](https://github.com/couchbaselabs/Checkers-iOS) from github.  Unfortunately it is a private repo at the time of this writing.  Contact [wacarter](https://github.com/wacarter) if you are interested in getting the source code.

It can be installed for the [iTunes Store](https://itunes.apple.com/us/app/id698034787), however that version is only able to connect to the non-public production server.

# Configure checkerlution

Edit SERVER_URL in gamecontroller.go to point the Sync Gateway you want to test against.

# Run checkerlution

```
$ cd $GOPATH/github.com/tleyden/checkerlution/main
$ go run main.go
```

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

