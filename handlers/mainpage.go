package handlers

import (
	"net/http"
	"text/template"
)

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
	data := struct{}{}

	// Render the template with the provided data
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "cannot render index.html", http.StatusInternalServerError)
	}
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/mainpage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
