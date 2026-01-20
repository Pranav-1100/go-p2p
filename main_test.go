package main

import (
	"net/http"
	"testing"
	"time"
)

func TestHealthEndpoint(t *testing.T) {
	go startHealthServer()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestExitCommand(t *testing.T) {
	if !isExitCommand("exit") {
		t.Fail()
	}
	if isExitCommand("hello") {
		t.Fail()
	}
}