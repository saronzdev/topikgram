package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"topikgram/api/internal/migrations"
	"topikgram/api/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	var direction = flag.String("dir", "up", "migration direction: up or down")
	flag.Parse()

	godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	ctx := context.Background()
	db, err := store.NewPool(ctx, connStr)
	if err != nil {
		log.Fatalf("Database failed: %v", err)
	}
	defer db.Close()

	m := migrations.New(db)

	switch *direction {
	case "up":
		if res := m.Up(ctx); res != nil {
			log.Fatal(res)
		} else {
			fmt.Println("Migration: Nothing to do")
		}
	default:
		log.Fatal("use: go run cmd/migrate/main.go -dir=up")
	}
}
