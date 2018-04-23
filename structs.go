package main

import (
	"net"
	"sync"
	"time"
)

type Nexus struct {
	Players map[string]Player
	Lock    *sync.Mutex
}

type Games struct {
	Games map[string]Game
	Lock  *sync.Mutex
}

type ConnMap struct {
	Conns map[string]ConnArray
	Lock  *sync.Mutex
}

type Game struct {
	Type, Map string
	Closed    bool
	Expires   time.Time
	Players   map[string]PlayerCell
}

type Player struct {
	Username, ID, Level, CurrentGame, Class string
	Created                                 time.Time
	Conn                                    net.Conn
	Friends                                 []Player
}

type PlayerCell struct {
	Name, Level, Class string
	K, D               string
	Conn               net.Conn
	Actions            net.Conn
}
