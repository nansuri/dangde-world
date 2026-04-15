package main

import (
	"log"

	"dangde-world/backend/internal/infrastructure/persistence"
	httpserver "dangde-world/backend/internal/interfaces/http"
	"dangde-world/backend/internal/shared/config"
)

func main() {
	cfg := config.Load()

	db, err := persistence.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := persistence.Seed(db); err != nil {
		log.Fatalf("database seed failed: %v", err)
	}

	router := httpserver.NewRouter(db, cfg)

	log.Printf("DangDe! World backend listening on %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
