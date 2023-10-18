package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	auth "forum/middleware"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
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
	// Generate a session token (you should implement this)
	sessionToken, err := auth.GenerateSessionToken(username)
	if err != nil {
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}
	auth.SetSessionCookie(w, sessionToken)
	http.Redirect(w, r, "/main-page", http.StatusSeeOther)
	fmt.Println(userID)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error rendering login page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Please fill in username or password", http.StatusBadRequest)
		return // Return to exit the function
	}
	// Retrieve the hashed password from the database for the given username.
	storedHashedPassword, err := auth.GetHashedPassword(db, username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use bcrypt.CompareHashAndPassword to check if the provided password matches the stored hashed password.
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	// http.Redirect(w, r, "/profile", http.StatusSeeOther)
	http.Redirect(w, r, "/main-page", http.StatusSeeOther)

}
