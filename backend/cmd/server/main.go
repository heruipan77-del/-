package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"smartops/backend/internal/api"
)

func main() {
	mux := http.NewServeMux()
	apiServer := &api.Server{}
	apiServer.Register(mux)

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = filepath.Join("..", "web", "dist")
	}
	mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

	addr := ":8080"
	log.Printf("smartops server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
