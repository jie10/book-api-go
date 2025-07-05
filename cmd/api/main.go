package main

import (
	"fmt"
	"log"
	"net/http"
)

func book(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Book API!!"))
}

func main() {
	port := ":4000"

	mux := http.NewServeMux()
	mux.HandleFunc("/", book)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Println("Server is running on port:", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting the server")
	}
}
