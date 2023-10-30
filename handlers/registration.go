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
	"golang.org/x/crypto/bcrypt"
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
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
	userID, err := database.GetIDBYusername(db, username)
	if err != nil {
		http.Error(w, "Failed to get user id", http.StatusInternalServerError)
	}
	fmt.Println(userID)
	sessionToken := uuid.New().String()

	// Set the session token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		// You can add more settings here like Expires, Secure, HttpOnly, etc.
	})

	// http.Redirect(w, r, "/profile", http.StatusSeeOther)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session token by setting the cookie to expire.
	http.SetCookie(w, &http.Cookie{
		Name:   "sessionToken",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expiration time in seconds (0 or negative value)
	})

	// Redirect the user to the login page or any other desired page.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
