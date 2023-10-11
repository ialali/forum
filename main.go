package main

import (
	"forum/database"
	"log"
)

func main() {
	// Define the path to your SQLite database file
	dbPath := "/Users/ibrahim/01Founders/forum/database/database.db"

	// Open a connection to the database
	db, err := database.OpenDatabase(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the schema and create tables
	err = database.InitializeSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	// You can now use the 'db' connection to perform database operations.
}
