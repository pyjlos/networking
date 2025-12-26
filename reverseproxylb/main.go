package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var backendId string // don't do this in prod - use a struct or something

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	backendId = os.Getenv("BACKEND_ID")
	if backendId == "" {
		backendId = "unknown-backend"
	}

	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", slowHandler)
	http.HandleFunc("/crash", crashHandler)
	http.HandleFunc("/health", healthHandler)

	log.Printf("Starting server on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Backend-ID", backendId)
	fmt.Fprintln(w, "Hello from fast endpoint!")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Backend-ID", backendId)
	fmt.Println("Received request at /slow, simulating delay...")
	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "Hello from slow endpoint")
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Backend-ID", backendId)
	// panic("Simulated server crash")
	os.Exit(1) // net/http doesn't crash by default -> it recovers -> force it to stop
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Backend-ID", backendId)
	fmt.Fprintln(w, "OK")
}
