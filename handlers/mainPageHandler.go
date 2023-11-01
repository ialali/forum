package handlers

import (
	"database/sql"
	"forum/database"
	"net/http"
	"text/template"
)

func ShowPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	posts, err := database.GetPosts(db)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
	}
	data := struct {
		Posts []database.Post
	}{
		Posts: posts,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error Parsing index.html", http.StatusInternalServerError)

	}
	tmpl.Execute(w, data)

}
