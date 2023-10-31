package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	auth "forum/middleware"
	"net/http"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	tmpl.Execute(w, nil)
}

func RegisterSubmitHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return // Return to exit the function
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return // Return to exit the function
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Please fill in username or password", http.StatusBadRequest)
		return // Return to exit the function
	}
	if auth.IsDuplicateUser(db, username, email) {
		http.Error(w, "Username or email is already taken", http.StatusConflict)
		return
	}

	userID, err := database.RegisterUser(db, username, email, password)
	if err != nil {
		http.Error(w, "Registration Failure", http.StatusInternalServerError)
		return // Return to exit the function
	}
	sessionToken := uuid.New().String()

	// Set the session token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		// You can add more settings here like Expires, Secure, HttpOnly, etc.
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
	fmt.Println(userID)

}
