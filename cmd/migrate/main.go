package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
)

func main() {
	// 1. Validate CLI arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	// 2. Establish connection to the SQLite database
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// 3. Create a driver instance for golang-migrate to interact with SQLite
	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")

	if err != nil {

		log.Fatal(err)

	}

	m, err := migrate.NewWithInstance("file", fSrc, "sqlite3", instance)

	if err != nil {

		log.Fatal(err)

	}

	// 5. Execute based on user input
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Printf("✅ Success: Migrated UP")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Printf("🔄 Success: Migrated DOWN")

	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'")
	}

}
