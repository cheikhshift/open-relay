package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func MatchNex(cmd string) {

	var game Game
	var counter time.Time

	for {
		sinc := time.Since(counter)
		seconds := sinc.Seconds()
		if seconds > 2 {
			game = GetGame(cmd)
			stattable := GenerateStats(game)

			duration := time.Since(game.Expires)
			durm := duration.Minutes()

			if durm > Maxminutes {
				ExitGame(cmd, game)
				break
			}

			if durm < Maxminutes {
				PushStat(duration, stattable, game, cmd)
			}

			counter = time.Now()
		}
	}

	go PrepareClean(game)

}

func MatchNexTDA(cmd string) {

	var game Game
	var counter time.Time

	for {
		sinc := time.Since(counter)
		seconds := sinc.Seconds()
		if seconds > 2 {
			game = GetGameTDA(cmd)
			stattable := GenerateStatsTDA(game)

			duration := time.Since(game.Expires)
			durm := duration.Minutes()

			if durm > Maxminutes {
				ExitGameTDA(cmd, game)
				break
			}

			if durm < Maxminutes {
				PushStatTDA(duration, stattable, game, cmd)
			}

			counter = time.Now()
		}
	}

	go PrepareClean(game)

}

func PrepareClean(game Game) {
	for playerid, _ := range game.Players {
		go CleanRelay(playerid)
	}
}

func CleanRelay(playerid string) {

	CM.Lock.Lock()
	conns, ok := CM.Conns[playerid]

	if !ok {
		CM.Lock.Unlock()
		return
	}

	for _, conn := range conns {
		if conn != nil {
			conn.Close()
		}
	}
	delete(CM.Conns, playerid)
	CM.Lock.Unlock()

}

func PushStat(duration time.Duration, stattable []string, game Game, gameid string) {
	minleft := Maxminutes - duration.Minutes()

	timefm := fmt.Sprintf("%.f minutes left.", minleft)
	joinstr := strings.Join(stattable, delim)
	joinstr = fmt.Sprintf("5%%%s%%%s%%\n", timefm, joinstr)
	strbytes := []byte(joinstr)
	for key, player := range game.Players {
		if player.Conn != nil {
			player.Conn.Write(strbytes)
		} else {
			Matches.Lock.Lock()
			delete(game.Players, key)
			Matches.Games[gameid] = game
			Matches.Lock.Unlock()
		}
	}
}

func PushStatTDA(duration time.Duration, stattable []string, game Game, gameid string) {
	minleft := Maxminutes - duration.Minutes()

	timefm := fmt.Sprintf("%.f minutes left.", minleft)
	joinstr := strings.Join(stattable, delim)
	joinstr = fmt.Sprintf("5%%%s%%%s%%\n", timefm, joinstr)
	strbytes := []byte(joinstr)
	for key, player := range game.Players {
		if player.Conn != nil {
			player.Conn.Write(strbytes)
		} else {
			MatchesTDA.Lock.Lock()
			delete(game.Players, key)
			MatchesTDA.Games[gameid] = game
			MatchesTDA.Lock.Unlock()
		}
	}
}

func ExitGame(gameid string, game Game) {

	Matches.Lock.Lock()
	delete(Matches.Games, gameid)
	Matches.Lock.Unlock()

	subset := "10%0"
	subsetbyte := []byte(subset)
	for _, player := range game.Players {
		if player.Conn != nil {
			player.Conn.Write(subsetbyte)
		}
	}
}

func ExitGameTDA(gameid string, game Game) {

	MatchesTDA.Lock.Lock()
	delete(MatchesTDA.Games, gameid)
	MatchesTDA.Lock.Unlock()

	subset := "10%0"
	var L, D int
	for _, player := range game.Players {

		dth, _ := strconv.Atoi(player.D)
		i, _ := strconv.Atoi(player.K)
		i -= dth
		if player.Team == "L" {
			L += i
		} else {
			D += i
		}
	}
	if L > D {
		subset = fmt.Sprintf("%s%%L", subset)
	} else {
		subset = fmt.Sprintf("%s%%D", subset)
	}

	subsetbyte := []byte(subset)
	for _, player := range game.Players {
		if player.Conn != nil {
			player.Conn.Write(subsetbyte)
		}
	}
}

func GenerateStats(game Game) (stattable []string) {

	for _, player := range game.Players {
		if player.Name != "" {
			subset := fmt.Sprintf("%s/%s/%s/%s/%s", player.Name, player.K,
				player.D, player.Level, player.Class)

			stattable = append(stattable, subset)
		}
	}
	return
}

func GenerateStatsTDA(game Game) (stattable []string) {

	for _, player := range game.Players {
		if player.Name != "" {
			subset := fmt.Sprintf("%s/%s/%s/%s/%s/%s", player.Name, player.K,
				player.D, player.Level, player.Class, player.Team)

			stattable = append(stattable, subset)
		}
	}
	return
}

func GetGame(gameid string) (game Game) {
	Matches.Lock.Lock()
	game = Matches.Games[gameid]
	Matches.Lock.Unlock()
	return
}

func GetGameTDA(gameid string) (game Game) {
	MatchesTDA.Lock.Lock()
	game = MatchesTDA.Games[gameid]
	MatchesTDA.Lock.Unlock()
	return
}
