package main

import (
	"net"
	"testing"
	"time"
)

func TestStartListener(t *testing.T) {
	done := make(chan net.Conn)
	errChan := make(chan error)

	go func() {
		conn, err := startListener(":9999")
		if err != nil {
			errChan <- err
			return
		}
		done <- conn
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Fatalf("Failed to dial listener: %v", err)
	}
	defer conn.Close()

	select {
	case listenerConn := <-done:
		if listenerConn == nil {
			t.Fatal("Listener returned nil connection")
		}
		defer listenerConn.Close()
	case err := <-errChan:
		t.Fatalf("Listener failed: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for listener to accept connection")
	}
}

func TestStartListenerInvalidAddress(t *testing.T) {
	_, err := startListener("invalid:address:format")
	if err == nil {
		t.Fatal("Expected error for invalid address, got nil")
	}
}

func TestDialPeer(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:9998")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}
	defer listener.Close()

	acceptDone := make(chan bool)
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
		acceptDone <- true
	}()

	conn, err := dialPeer("127.0.0.1:9998")
	if err != nil {
		t.Fatalf("dialPeer failed: %v", err)
	}

	if conn == nil {
		t.Fatal("dialPeer returned nil connection")
	}
	defer conn.Close()

	<-acceptDone
}

func TestDialPeerConnectionRefused(t *testing.T) {
	_, err := dialPeer("127.0.0.1:9997")
	if err == nil {
		t.Fatal("Expected error when connecting to closed port, got nil")
	}
}

func TestDialPeerInvalidAddress(t *testing.T) {
	_, err := dialPeer("not-a-valid-address")
	if err == nil {
		t.Fatal("Expected error for invalid address, got nil")
	}
}
