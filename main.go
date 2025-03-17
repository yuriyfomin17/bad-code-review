package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/yuriyfomin17/bad-code-review/config"
	"github.com/yuriyfomin17/bad-code-review/http_handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error happened while reading config", err)
	}
	httpHandler, err := http_handler.NewHttpServer(cfg.HttpClientTimeoutSeconds, cfg.NumOfWorkers, cfg.BatchNumOrdersIdsToProcess)
	if err != nil {
		log.Fatal("error happened while creating http server", err)
	}
	http.HandleFunc("/orders", httpHandler.OrderHandler)
	log.Printf("Starting HTTP server on %s", cfg.Port)
	err = http.ListenAndServe(cfg.Port, nil)
	if err != nil {
		log.Fatal("error happened while starting up the server", err)
	}
}
