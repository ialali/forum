package handlers

import (
	"database/sql"
	"forum/database"
	"log"
	"net/http"
	"text/template"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, "error parsing the template", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}

func AddPostSubmit(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if the user is authenticated
	userID, isAuthenticated := database.GetAuthenticatedUserID(r)
	if !isAuthenticated {
		// Redirect to the login page or show an error message.
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	category := r.FormValue("category")
	title := r.FormValue("title")
	content := r.FormValue("content")

	if category == "" || title == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusUnprocessableEntity)
		return
	}

	var categoryID int
	err := db.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID)
	if err != nil {
		if err != sql.ErrNoRows {
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			log.Fatal(err)
			return
		}

		// If the category doesn't exist, create it
		_, err = db.Exec("INSERT INTO categories (name) VALUES (?)", category)
		if err != nil {
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			log.Println(err)
			return
		}

		// Retrieve the newly created category ID
		err = db.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&categoryID)
		if err != nil {
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			log.Println(err)
			return
		}

		// Add the post to the database
		err = database.InsertPost(db, category, title, content, userID)
		if err != nil {
			log.Println("Error inserting post:", err)
			http.Redirect(w, r, "/error/500", http.StatusSeeOther)
			return
		}

		// Redirect the user to the home page after successfully adding the post
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
