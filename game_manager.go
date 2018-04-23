package main

import (
	"net"
	"fmt"
	"github.com/cheikhshift/gos/core"
	"time"
)

func AddPlayerToGame(cmd []string, conn net.Conn) {
	Matches.Lock.Lock()

	if match, exists := Matches.Games[cmd[1]]; exists {

		CheckifJoined(cmd[3], match)
		Sc.Lock.Lock()
		lpl := Sc.Players[cmd[3]]
		Sc.Lock.Unlock()
		match.Players[cmd[3]] = PlayerCell{Name: cmd[3], Level: lpl.Level, Class: lpl.Class, K: "0", D: "0", Conn: conn}
		Matches.Games[cmd[1]] = match
	}

	Matches.Lock.Unlock()
}

func CheckifJoined(playerid string, match Game) {
	_, exists := match.Players[playerid]

	if !exists {
		go CleanRelay(playerid)
	}
}

func RegisterStat(cmd []string) {
	Matches.Lock.Lock()
	game := Matches.Games[cmd[1]]

	playerc := game.Players[cmd[2]]
	playerc.K = cmd[3]
	playerc.D = cmd[4]
	game.Players[cmd[2]] = playerc
	Matches.Games[cmd[1]] = game
	Matches.Lock.Unlock()
}

func MatchMake(cmd []string) (str string) {

	var gameid string

	str, foundGame := FindGame(cmd)

	if !foundGame {
		gameid, str = AddNewGame(cmd)
		go MatchNex(gameid)
	}

	return
}

func FindGame(cmd []string) (response string, foundGame bool) {
	Matches.Lock.Lock()
	games := Matches.Games
	Matches.Lock.Unlock()

	for gameid, match := range games {
		playerlen := len(match.Players)

		if match.Type == cmd[1] && playerlen < MaxPlayers {
			response = fmt.Sprintf("%s%s%s%s%s", "1", delim, gameid, delim, match.Map)
			foundGame = true
		}
	}

	return
}

func AddNewGame(cmd []string) (key string, response string) {
	Matches.Lock.Lock()
	expire := time.Now()
	rint := random(0, 2)
	maps := []string{"Level8", "Level3", "Level4"}
	nmatch := Game{Type: cmd[1], Map: maps[rint], Players: make(map[string]PlayerCell), Expires: expire}
	key = core.NewLen(25)
	Matches.Games[key] = nmatch
	response = fmt.Sprintf("%s%s%s%s%s", "1", delim, key, delim, nmatch.Map)
	Matches.Lock.Unlock()
	return
}