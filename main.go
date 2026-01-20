package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	listenAddr := flag.String("listen", "", "Address to listen on")
	connectAddr := flag.String("connect", "", "Address to connect to")
	ciMode := flag.Bool("ci", false, "Run in CI mode")

	flag.Parse()

	if *ciMode {
		fmt.Println("Starting in CI Mode (Non-Interactive)...")
		go startHealthServer()
		select {}
	}

	if *listenAddr == "" && *connectAddr == "" {
		fmt.Println("Error: Specify -listen or -connect")
		os.Exit(1)
	}

	if *listenAddr != "" && *connectAddr != "" {
		fmt.Println("Error: Choose only one mode")
		os.Exit(1)
	}

	var conn net.Conn
	var err error

	if *listenAddr != "" {
		conn, err = startListener(*listenAddr)
	} else {
		conn, err = dialPeer(*connectAddr)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handleChat(conn)
}
