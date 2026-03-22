package main

import (
	"log"
	"net/http"

	"github.com/herdifirdausss/belajar-vibe-coding/internal/api"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/config"
	"github.com/herdifirdausss/belajar-vibe-coding/internal/db"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer database.Close()

	server := api.NewServer(database)
	mux := server.Routes()

	log.Printf("Starting server on port %s", cfg.ServerPort)
	err = http.ListenAndServe(":"+cfg.ServerPort, mux)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
