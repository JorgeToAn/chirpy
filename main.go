package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Started server at port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
