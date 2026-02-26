package main

import (
	"log"
	"os"
	"time"

	"github.com/crimsonn/zm_server/internal/database"
	"github.com/crimsonn/zm_server/internal/env"
	"github.com/crimsonn/zm_server/internal/server"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrationCommand()
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "generate_temp_jwt" {
		generateTempJWTCommand()
		return
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, env.GetEnvString("DB_MIGRATIONS_DIR", "./migrations")); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	startServer()
}

func generateTempJWTCommand() {

	secret := env.GetEnvString("FIVEM_BACKEND_SECRET", "")
	if secret == "" {
		log.Fatalf("FIVEM_BACKEND_SECRET is not configured")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "backend",
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("failed to generate temp jwt: %v", err)
	}
	log.Println("temp jwt:", tokenString)
}

func runMigrationCommand() {
	if len(os.Args) < 3 {
		log.Fatalf("missing migration action. usage: go run . migrate [up|down]")
	}

	action := os.Args[2]
	migrationsDir := env.GetEnvString("DB_MIGRATIONS_DIR", "./migrations")

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}
	defer db.Close()

	switch action {
	case "up":
		if err := database.RunMigrations(db, migrationsDir); err != nil {
			log.Fatalf("failed to run migrations up: %v", err)
		}
		log.Println("migrations up completed")
	case "down":
		if err := database.RollbackLastMigration(db, migrationsDir); err != nil {
			log.Fatalf("failed to run migration down: %v", err)
		}
		log.Println("migration down completed")
	default:
		log.Fatalf("invalid migration action %q. usage: go run . migrate [up|down]", action)
	}
}

func startServer() {
	cfg := server.Config{
		Host: env.GetEnvString("APP_HOST", "0.0.0.0"),
		Port: env.GetEnvNumber("APP_PORT", 8080),
	}

	srv := server.NewServer(cfg)
	done := make(chan bool)
	go func() {
		if err := srv.Start(done); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done
	log.Println("server shutdown complete")
	os.Exit(0)
}
