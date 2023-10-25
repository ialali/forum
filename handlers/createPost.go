package handlers

import (
	"database/sql"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Check for the presence of a session token in the request's cookies.
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		// No session token found, user is not authenticated.
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get the userID associated with the username.
	userID, isAuthenticated := GetUserIDForSessionToken(sessionToken.Value)
	if !isAuthenticated {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse the post data from the form.
	category := r.FormValue("category")
	title := r.FormValue("title")
	content := r.FormValue("content")

	// Insert the post into the database with the associated userID.
	_, err = db.Exec("INSERT INTO posts (category, title, content, user_id) VALUES (?, ?, ?, ?)", category, title, content, userID)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Redirect to the main page or another appropriate location.
	http.Redirect(w, r, "/main-page", http.StatusSeeOther)
}
