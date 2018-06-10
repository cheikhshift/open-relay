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

	hostname := flag.String("hostname", "", "Hostname the server should listen on.")
	port := flag.Int("port", 3333, "Port number to listen on.")
	ip := GetLocalIP()
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *hostname == "" {
		hostname = &ip
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	err := LoadConfig(&Sc)

	if err != nil {
		Sc = &Nexus{Players: make(map[string]Player)}
		fmt.Println(err.Error())

	}

	go SaveOnExit(stop)

	Sc.Lock = new(sync.Mutex)

	serverAdr := fmt.Sprintf("%s:%v", *hostname, *port)

	l, err := net.Listen(CONN_TYPE, serverAdr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Welcome to Warrior project multiplayer server. Now open your game, tap on `configure online` and enter the following information.")
	fmt.Println("Have fun. :)")
	fmt.Println("IP Address :  ", ip)
	fmt.Println("Port : ", *port)

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

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
