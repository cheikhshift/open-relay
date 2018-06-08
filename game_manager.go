package main

import (
	"crypto/rand"
	"fmt"
	"net"
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

func AddPlayerToGameTDA(cmd []string, conn net.Conn) {
	MatchesTDA.Lock.Lock()

	if match, exists := MatchesTDA.Games[cmd[1]]; exists {

		CheckifJoined(cmd[3], match)
		Sc.Lock.Lock()
		lpl := Sc.Players[cmd[3]]
		Sc.Lock.Unlock()
		var tm string
		lt, dk := CountTeams(match.Players)

		if lt <= dk {
			tm = "L"
		} else {
			tm = "D"
		}

		match.Players[cmd[3]] = PlayerCell{Name: cmd[3], Level: lpl.Level, Class: lpl.Class, K: "0", D: "0", Team: tm, Conn: conn}
		MatchesTDA.Games[cmd[1]] = match
	}

	MatchesTDA.Lock.Unlock()
}

func CountTeams(players map[string]PlayerCell) (lt int, dk int) {

	for _, dta := range players {
		if dta.Team == "L" {
			lt++
		} else {
			dk++
		}
	}

	return
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

func RegisterStatTDA(cmd []string) {
	MatchesTDA.Lock.Lock()
	game := MatchesTDA.Games[cmd[1]]

	playerc := game.Players[cmd[2]]
	playerc.K = cmd[3]
	playerc.D = cmd[4]
	game.Players[cmd[2]] = playerc
	MatchesTDA.Games[cmd[1]] = game
	MatchesTDA.Lock.Unlock()
}

func MatchMakeTDA(cmd []string) (str string) {
	var gameid string

	str, foundGame := FindGameTDA(cmd)

	if !foundGame {
		gameid, str = AddNewGameTDA(cmd)
		go MatchNexTDA(gameid)
	}

	return
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

func FindGameTDA(cmd []string) (response string, foundGame bool) {
	MatchesTDA.Lock.Lock()
	games := MatchesTDA.Games
	MatchesTDA.Lock.Unlock()

	for gameid, match := range games {
		playerlen := len(match.Players)

		if match.Type == cmd[1] && playerlen < MaxPlayers {
			response = fmt.Sprintf("%s%s%s%s%s", "1", delim, gameid, delim, match.Map)
			foundGame = true
		}
	}

	return
}

func NewLen(length int) string {

	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	if length == 0 {
		return ""
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		panic("uniuri: wrong charset length for NewLenChars")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("uniuri: error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

func AddNewGame(cmd []string) (key string, response string) {
	Matches.Lock.Lock()
	expire := time.Now()
	rint := random(0, 2)
	maps := []string{"Level2", "Level8", "Level3", "Level4"}
	nmatch := Game{Type: cmd[1], Map: maps[rint], Players: make(map[string]PlayerCell), Expires: expire}
	key = NewLen(25)
	Matches.Games[key] = nmatch
	response = fmt.Sprintf("%s%s%s%s%s", "1", delim, key, delim, nmatch.Map)
	Matches.Lock.Unlock()
	return
}

func AddNewGameTDA(cmd []string) (key string, response string) {
	MatchesTDA.Lock.Lock()
	expire := time.Now()
	rint := random(0, 2)
	maps := []string{"Level2", "Level8", "Level3", "Level4"}
	nmatch := Game{Type: cmd[1], Map: maps[rint], Players: make(map[string]PlayerCell), Expires: expire}
	key = NewLen(25)
	MatchesTDA.Games[key] = nmatch
	response = fmt.Sprintf("%s%s%s%s%s", "1", delim, key, delim, nmatch.Map)
	MatchesTDA.Lock.Unlock()
	return
}
