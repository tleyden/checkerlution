[![Build Status](https://drone.io/github.com/tleyden/checkerlution/status.png)](https://drone.io/github.com/tleyden/checkerlution/latest)

A checkers bot implementation for [checkers-bot](https://github.com/tleyden/checkers-bot) which uses a neural network to do its thinking.  The same essential approach is used as described in [Blondie24: Playing at the edge of AI](http://www.amazon.com/Blondie24-Playing-Kaufmann-Artificial-Intelligence/dp/1558607838).  

The underlying neural network library code is found in [neurgo](https://github.com/tleyden/neurgo) and [neurvolve](https://github.com/tleyden/neurvolve).

# Big Picture

![architecture png](http://cl.ly/image/1V1D393S0A45/Screen%20Shot%202013-10-13%20at%2010.53.01%20AM.png)

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

# Configure checkerlution

There are command line options to tell it which server to use.  See main/run.sh for examples.

# Run checkerlution

```
$ cd $GOPATH/github.com/tleyden/checkerlution/main
$ go run main.go
```
