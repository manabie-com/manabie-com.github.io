package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

var addr = os.Getenv("ADDR")

func main() {
	// not used...
	Serve()
}

var server http.Server

func Serve() {
	mux := http.NewServeMux()

	mux.Handle("/a", http.HandlerFunc(HandlerA))
	mux.Handle("/b", http.HandlerFunc(HandlerB))
	mux.Handle("/kill", http.HandlerFunc(KillHandler))
	server = http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// block until server Shutdown
	log.Println("Starting server at", addr)
	log.Println(server.ListenAndServe())
}

func HandlerA(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A")
	w.WriteHeader(http.StatusOK)
}

func HandlerB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("B")
	fmt.Println("has")
	fmt.Println("2 extra lines of code")
	w.WriteHeader(http.StatusOK)
}

func KillHandler(w http.ResponseWriter, r *http.Request) {
	server.Shutdown(context.Background())
	w.WriteHeader(http.StatusOK)
}
