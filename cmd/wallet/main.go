package main

import (
	"fmt"
	"log"
	"net/http"

	"wallet-service/internal/config"
	"wallet-service/internal/handler"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/createaddress", handler.CreateAddress(cfg))
	mux.HandleFunc("/api/v1/validateaddress", handler.ValidateAddress(cfg))
	mux.HandleFunc("/api/v1/tx", handler.SignTx(cfg))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting wallet service on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
