package handlers

import (
	"forum/database"
	"net/http"
	"text/template"
)

// var post database.Post

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	// Load your main page template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "cannot parse index.html", http.StatusInternalServerError)
		return
	}

	// You can fetch data from your database and pass it to the template
	// For example, you can retrieve user-specific data to display on the main page.

	// In this example, we're not passing any data to the template.
	data := struct {
		Posts []database.Post
	}{}

	// Render the template with the provided data
	tmpl.Execute(w, data)
}


func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
