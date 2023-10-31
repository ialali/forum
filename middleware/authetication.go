package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	sessions = make(map[string]int)
	mu       sync.Mutex
)

func IsDuplicateUser(db *sql.DB, username, email string) bool {
	// Perform a database query to check if the username or email already exists in the database.
	// Return true if duplicate, false if not.
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = ? OR email = ?"
	row := db.QueryRow(query, username, email)
	if err := row.Scan(&count); err != nil {
		// Handle the error, e.g., log it or return false.
		return false
	}

	return count > 0
}
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
func GetHashedPassword(db *sql.DB, username string) (string, error) {
	// Query the database to get the hashed password for the provided username.
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		// Handle errors, e.g., username not found in the database.
		if err == sql.ErrNoRows {
			return "", errors.New("User not found")
		}
		return "", err
	}

	return hashedPassword, nil
}

func SetSessionCookie(w http.ResponseWriter, userID int) {
	sessionID := uuid.New().String()
	// Set the session token as a cookie
	cookie := http.Cookie{
		Name:    "session",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
		// You can set other cookie properties such as expiration, path, secure, HttpOnly, etc.
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	sessions[sessionID] = userID
}

// func GenerateSessionToken(username string) (string, error) {
// 	// You can use a combination of username, current timestamp, and some secret key
// 	// to generate a unique session token.

// 	// In this example, we concatenate the username and current timestamp.
// 	tokenData := username + time.Now().Format("20060102150405")

// 	// You can hash or encrypt the token data for added security.
// 	// For example, you can use a package like crypto/sha256 or crypto/md5.
// 	hashedTokenData := sha256.Sum256([]byte(tokenData))
// 	token := hex.EncodeToString(hashedTokenData[:])

//		return token, nil
//	}
func GetAuthenticatedUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session-name")
	if err != nil {
		return 0, false
	}
	userID, ok := sessions[cookie.Value]
	return userID, ok
}
func IsAuthenticated(r *http.Request) bool {
	// Check if the user is authenticated by looking for a session token.
	cookie, err := r.Cookie("session_token")
	if err != nil {
		// No session token found, the user is not authenticated.
		log.Println("No session token found.")
		return false
	}

	// Retrieve the session token from the cookie.
	sessionToken := cookie.Value
	log.Println("Retrieved session token:", sessionToken)

	// Look up the user's ID associated with the session token.
	mu.Lock()
	defer mu.Unlock()
	_, ok := sessions[sessionToken]

	if ok {
		log.Println("User is authenticated.")
	} else {
		log.Println("User is not authenticated.")
	}

	// If the session token is found in the sessions map, the user is authenticated.
	return ok
}
