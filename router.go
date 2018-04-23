package main

import (
	"fmt"
	"github.com/cheikhshift/gos/core"
	"net"
	"strings"
	"time"
)

func NBox(command string, conn net.Conn, playerID string) (retval string, err error) {

	cmd := strings.Split(command, delim)
	var response []byte
	var str string

	switch cmd[0] {
	case "register":
		str, retval = RegisterUser(cmd, conn)

	case "register-join":

		AddPlayerToGame(cmd, conn)
		str = "2%0"

	case "register-metric":

		RegisterStat(cmd)
	case "register-action":
		Matches.Lock.Lock()
		_, ok := Matches.Games[cmd[1]]
		Matches.Lock.Unlock()

		if ok {

			go RelayConn(conn, cmd[2])
		}
		str = "OK%"
		retval = "disto"

	case "register-read":

		CM.Lock.Lock()
		CM.Conns[cmd[2]] = append(CM.Conns[cmd[2]], conn)
		CM.Lock.Unlock()
		retval = "noread"
		// go Movbox(cmd, conn)

	case "register-findffa":

		str = MatchMake(cmd)
	default:
	}

	response = []byte(str)
	conn.Write(response)

	return
}

// Start Nbox supporting funcs
func RegisterUser(cmd []string, conn net.Conn) (str string, retval string) {
	Sc.Lock.Lock()
	hasdelim := strings.Contains(cmd[1], "/") || strings.Contains(cmd[1], "%") || len(cmd[1]) > 13
	if result, ext := Sc.Players[cmd[1]]; (result.ID == cmd[4] || !ext) && !hasdelim {
		str = fmt.Sprintf("%s%s%s", "0", delim, "0")
		player := Player{Username: cmd[1], Level: cmd[2], Friends: []Player{}, ID: cmd[4], Class: cmd[3], Conn: conn}

		retval = cmd[1]
		player.CurrentGame = ""
		Sc.Players[cmd[1]] = player
		fmt.Println(player, "<----")
	} else {
		str = fmt.Sprintf("%s%s%s", "0", delim, "1")
	}
	Sc.Lock.Unlock()

	return
}

func AddPlayerToGame(cmd []string, conn net.Conn) {
	Matches.Lock.Lock()
	if match, exists := Matches.Games[cmd[1]]; exists {
		Sc.Lock.Lock()
		lpl := Sc.Players[cmd[3]]
		Sc.Lock.Unlock()
		match.Players[cmd[3]] = PlayerCell{Name: cmd[3], Level: lpl.Level, Class: lpl.Class, K: "0", D: "0", Conn: conn}
		Matches.Games[cmd[1]] = match
	}

	Matches.Lock.Unlock()
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

	str,foundGame := FindGame(cmd)
	
	if !foundGame {
		gameid, str = AddNewGame(cmd)
		go MatchNex(gameid)
	}

	return
}

func FindGame(cmd []string) (response string, foundGame bool){
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
