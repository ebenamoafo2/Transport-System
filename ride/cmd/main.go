package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ebenamoafo2/transport/ride/configs"
)

func main() {
	configPath := flag.String("config", "../configs/config.yml", "path to YAML configuration")
	flag.Parse()
	cfg, err := configs.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	//Mask secrets before logging
	safe := *cfg
	safe.Database.Password = "<redacted>"
	log.Printf("Loaded config: %+v", safe)

	mux := http.NewServeMux()
	//Root endpoint: a friendly hello
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "Hello from the ride service!")
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			name = "stranger"
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"message": fmt.Sprintf("Hello, %s", name),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("failed to decode response: %v", err)
		}
	})

	// /slow Simulates slow work, e.g. a long DB query, so we can test shutdown behavior.
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("slow request started")
		time.Sleep(8 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "finished slow work"}`)
		log.Println("slow request finished")
	})

	//Liveness probe: reports the process is alive.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ok")
	})

	// Create the HTTP server using values from the configuration.
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSec) * time.Second,
	}
	log.Printf("Starting server on port %d", cfg.Server.Port)

	//Listen for shutdown signals
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		// Start the HTTP server.
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("listen failure on %s: %w", server.Addr, err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-sigCtx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown timed out, forcing close: %v", err)
			_ = server.Close()
		}

	case err := <-errCh:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}
		log.Println("server exited")
		return
	}

}
