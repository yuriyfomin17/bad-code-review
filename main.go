package main

import (
	"bad-code-review/config"
	"bad-code-review/http_handler"
	"log"
	"net/http"
)

func main() {
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
