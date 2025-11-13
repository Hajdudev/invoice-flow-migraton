package main

import (
	"flag"
	"log"
	"os"

	database "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql"
	"github.com/Hajdudev/invoice-flow/internal/adapters/postgresql/migrations"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "up", "Migration action: up, down, down-all")
	flag.Parse()

	log.Printf("Starting database migration (%s)...", action)

	db := database.New()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	var err error
	switch action {
	case "up":
		err = db.MigrateFS(migrations.FS, ".")
	case "down":
		err = db.MigrateDownFS(migrations.FS, ".")
	case "down-all":
		err = db.MigrateDownAllFS(migrations.FS, ".")
	default:
		log.Fatalf("Unknown action: %s. Use 'up', 'down', or 'down-all'", action)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
		os.Exit(1)
	}

	log.Println("Migration completed successfully!")
}
