// main.go
package main

import (
	"log"
	"net/http"

	"github.com/ShawnMa123/sub-node-cvt/internal/handler"
)

func main() {
	// API handler
	http.HandleFunc("/sub", handler.SubscriptionHandler)

	// Static file server for the frontend
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	port := "8080"
	log.Printf("Server starting on http://localhost:%s", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}