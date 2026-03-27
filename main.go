package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricsGet)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Started server at port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
