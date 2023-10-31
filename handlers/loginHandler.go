package handlers

import (
	"database/sql"
	"fmt"
	"forum/database"
	auth "forum/middleware"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	sessions     = make(map[string]int)
	userSessions = make(map[int]string)
	mu           sync.Mutex
)

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
	log.Println("Generated session token:", sessionToken)

	// Store the session token in your server (assuming you have sessions and userSessions maps)
	mu.Lock()
	defer mu.Unlock()
	// Store the session ID and user ID in their respective maps
	userSessions[userID] = sessionToken
	sessions[sessionToken] = userID

	// Set the session token as a cookie in the response
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour), // Adjust the expiration time
		HttpOnly: true,
		// You can add more settings here like Secure, SameSite, etc.
	})

	// Redirect the user after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println(sessions)
}
