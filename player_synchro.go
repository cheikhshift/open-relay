package main

import (
	"net"
	"time"
)

func RelayConn(conn net.Conn, playerid string) {
	buf := make([]byte, 512)
	// Read the incoming connection into the buffer.
	var counter time.Time
	for {
		sinc := time.Since(counter)
		seconds := sinc.Nanoseconds()
		if seconds > Latency {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}

			if n > 0 {
				CM.Lock.Lock()
				conns := CM.Conns[playerid]
				//fmt.Println(playerid)
				for _, conn := range conns {
					if conn != nil {
						//error muted conn.Write returns error
						conn.Write(buf[:n])
					}
				}
				CM.Lock.Unlock()
			}
			counter = time.Now()
		}
	}
	conn.Close()
}
