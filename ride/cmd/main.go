package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ebenamoafo2/transport/ride/configs"
	"github.com/ebenamoafo2/transport/ride/internal/adapters/repository"
	"github.com/ebenamoafo2/transport/ride/internal/httpserver"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()
	cfg, err := configs.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	safe := *cfg
	safe.Database.Password = "<redacted>"
	log.Printf("Loaded config: %+v", safe)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	// Pool configuration (database/sql pool lives inside *sql.DB)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)
	db.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTime) * time.Second)

	assignmentRepo := repository.NewSQLAssignmentRepository(db)
	if err := httpserver.Run(cfg.Server, assignmentRepo); err != nil {
		log.Fatalf("application run failed: %v", err)
	}
}
