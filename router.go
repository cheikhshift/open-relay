package main

import (
	"fmt"
	"net"
	"strings"
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

