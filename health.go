package main

import (
	"encoding/json"
	"net/http"
)

func startHealthServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "UP",
		})
	})
	http.ListenAndServe(":8081", nil)
}
