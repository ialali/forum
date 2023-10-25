package main

import (
	"fmt"
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlers.IndexPageHandler)

	http.HandleFunc("/register", handlers.IndexPageHandler)
	http.HandleFunc("/registerauth", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterSubmitHandler(w, r, db)
	})
	http.HandleFunc("/main-page", handlers.MainPageHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/loginauth", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginSubmitHandler(w, r, db)
	})
	http.HandleFunc("/logout", handlers.LogoutHandler)

	fmt.Println("server started on http://localhost:1212")
	http.ListenAndServe(":1212", nil)

	// You can now use the 'db' connection to perform database operations.
}
