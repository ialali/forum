package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	"log"
	"net/http"
	"text/template"
)

func ShowPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	posts, err := database.GetPosts(db)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
	}

	for i, post := range posts {
		comments, err := database.GetCommentsForPost(db, post.ID)
		if err != nil {
			// If there's an error fetching comments, return a 500 error and log the error
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		// Assign the comments to the post
		posts[i].Comments = comments
	}

	data := struct {
		Posts []database.Post
	}{
		Posts: posts,
	}
	fmt.Println(posts)
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error Parsing index.html", http.StatusInternalServerError)

	}
	tmpl.Execute(w, data)

}
