package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"api-gateway/clients"
	"api-gateway/handlers"
	"api-gateway/models"
)

func main() {
	config, err := loadConfig("config/routes.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	clientFactory := clients.NewGRPCClientFactory(config)
	grpcHandler := handlers.NewGRPCHandler(clientFactory, config.Routes)

	for endpoint := range config.Routes {
		http.HandleFunc(endpoint, grpcHandler.HandleRequest)
	}

	log.Println("API Gateway running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfig(filePath string) (*models.Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config models.Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	return &config, nil
}
