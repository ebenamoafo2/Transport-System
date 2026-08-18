package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ebenamoafo2/transport/ride/configs"
)

func main() {
	configPath := flag.String("config", "configs/config.yml", "path to YAML configuration")
	flag.Parse()
	cfg, err := configs.LoadConfig(*configPath)
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
		ReadTimeout: time.Duration(cfg.Server.ReadTimeoutSec),
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSec),
		
	}
	log.Printf("Loaded config: %+v", safe)
	log.Printf("Starting server on port %d", cfg.Server.Port)

	// Start the HTTP server.
	err = server.ListenAndServe()
	log.Fatal(err)
}
