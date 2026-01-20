package main

import (
	"fmt"
	"net"
)

func startListener(address string) (net.Conn, error) {
	fmt.Printf("⏳ Listening on %s...\n", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to start listener: %w", err)
	}
	conn, err := listener.Accept()
	if err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to accept: %w", err)
	}
	listener.Close()
	fmt.Printf("✓ Peer connected from %s\n", conn.RemoteAddr())
	return conn, nil
}

func dialPeer(address string) (net.Conn, error) {
	fmt.Printf("⏳ Connecting to %s...\n", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	fmt.Printf("✓ Connected to %s\n", conn.RemoteAddr())
	return conn, nil
}
