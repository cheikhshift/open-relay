package main

import "sync"

var MaxPlayers int = 6

var Latency int64 = 12000000 // in Nanosecond

const (
	CONN_TYPE = "tcp"
	delim     = "%"
	//Maximum game time
	Maxminutes float64 = 10
)

var (
	Sc         *Nexus
	Matches    *Games   = &Games{Games: make(map[string]Game), Lock: new(sync.Mutex)}
	CM         *ConnMap = &ConnMap{Conns: make(map[string]ConnArray), Lock: new(sync.Mutex)}
	MatchesTDA *Games   = &Games{Games: make(map[string]Game), Lock: new(sync.Mutex)}
)
