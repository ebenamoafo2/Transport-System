package main

import (
	"fmt"
	"log"
	"net/http"

	"gihub.com/ebenamoafo2/transport/ride/configs"
)

func main() {
	cfg, err := configs.LoadConfig("../configs/config.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	//Mask secrets before logging
	safe := *cfg
	safe.Database.Password = "<redacted>"
	log.Printf("Loaded config: %+v", safe)

	// Create the HTTP server using values from the configuration.
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Server.Port),
		IdleTimeout: cfg.Server.IdleTimeout,
	}
	log.Printf("Loaded config: %+v", safe)
	log.Printf("Starting server on port %d", cfg.Server.Port)

	// Start the HTTP server.
	err = server.ListenAndServe()
	log.Fatal(err)
}
