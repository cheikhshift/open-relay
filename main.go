package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {

	wordPtr := flag.String("hostname", "localhost", "Hostname the server should listen on.")
	numbPtr := flag.Int("port", 3333, "Port number to listen on.")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	// Listen for incoming connections.
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	err := LoadConfig(&Sc)

	if err != nil {
		Sc = &Nexus{Players: make(map[string]Player)}
		fmt.Println(err.Error())

	}

	Sc.Lock = new(sync.Mutex)

	serverAdr := fmt.Sprintf("%s:%v", *wordPtr, *numbPtr)

	l, err := net.Listen(CONN_TYPE, serverAdr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on ", serverAdr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}
