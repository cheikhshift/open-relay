package main

import "math/rand"
import "os"

func InArr(arr []Player, id string) bool {
	for _, player := range arr {
		if player.ID == id {
			return true
		}
	}
	return false
}

func random(min int, max int) int {
	return rand.Intn(max-min) + min
}

func SaveOnExit(stop chan os.Signal) {
	<-stop
	SaveConfig(&Sc)
	os.Exit(0)
}
