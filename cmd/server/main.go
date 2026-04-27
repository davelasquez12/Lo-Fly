package main

import (
	"log"
	"net/http"

	"LoflyBE/internal/config"
	"LoflyBE/internal/httpapi"
	"LoflyBE/internal/tracker/repository"
	"LoflyBE/internal/tracker/service"
)

func main() {
	cfg := config.Load()

	repo := repository.NewTrackerRepo(cfg.DataFile)
	trackerService := service.NewService(repo)

	server := httpapi.NewServer(trackerService, repo)
	log.Printf("flight price tracker backend listening on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, server.Routes()))
}
