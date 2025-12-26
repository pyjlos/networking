package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting simple http server")
	http.HandleFunc("/fast", fastHandler)
	http.HandleFunc("/slow", slowHandler)
	http.HandleFunc("/crash", crashHandler)
	http.ListenAndServe(":8080", nil)
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from fast endpoint!")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request at /slow, simulating delay...")
	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "Hello from slow endpoint")
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
	panic("Simulated server crash")
}
