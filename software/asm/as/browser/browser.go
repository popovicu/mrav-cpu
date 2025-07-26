package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
)

//go:embed as.html assembler.wasm wasm_exec.js
var staticFiles embed.FS

func main() {
	port := flag.Int("port", 9876, "http serving port")
	flag.Parse()

	handler := http.FileServerFS(staticFiles)
	// TODO: find out why the mux handler won't work
	// mux := http.NewServeMux()
	// mux.Handle("GET /as/{file}", handler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error running HTTP: %v\n", err)
	}
}
