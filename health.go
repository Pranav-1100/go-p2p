package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func startHealthServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{
			"status": "UP",
		}); err != nil {
			fmt.Printf("Error encoding health response: %v\n", err)
		}
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Health server error: %v\n", err)
	}
}