package main

import (
	"fmt"
	"net"
)

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 512)
	// Read the incoming connection into the buffer.
	var playerID string
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			Sc.Lock.Lock()
			player := Sc.Players[playerID]
			player.Conn = nil
			player.CurrentGame = ""
			Sc.Players[playerID] = player
			Sc.Lock.Unlock()
			break
		} else {
			msg := string(buf[:n])
			str, _ := NBox(msg, conn, playerID)
			if len(str) > 0 {
				if str == "disto" || str == "noread" {
					playerID = str
					break
				} else {
					playerID = str
				}
			}
		}
		// Send a response back to person contacting us.
		//conn.Write([]byte("Message received."))
	}

	// Close the connection when you're done with it.
	if playerID != "disto" && playerID != "noread" {
		conn.Close()
	}
}
