package database

import (
	"database/sql"
	auth "forum/middleware"
	"time"
)

func RegisterUser(db *sql.DB, username, email, password string) (int64, error) {
	// Hash the password before inserting it into the database (assuming you've set up bcrypt).
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	// Get the current registration date.
	registrationDate := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec(`
        INSERT INTO users (username, email, password, registration_date)
        VALUES (?, ?, ?, ?);
    `, &username, &email, &hashedPassword, registrationDate)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow(`
		SELECT id, username, email, password, registration_date
		FROM users
		WHERE id = ?;
	`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RegistrationDate)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id, email, password, registration_date FROM users WHERE email = ?`, email).Scan(&user.ID, &user.Email, &user.Username, &user.RegistrationDate)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
func GetIDBYusername(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&user.ID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// CreatePost inserts a new post into the database and returns the post ID.
func InsertPost(db *sql.DB, username, title, content, category string) (int64, error) {
	// Get the current timestamp for the post
	postTime := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec(`
        INSERT INTO posts (username, title, content, category, post_time)
        VALUES (?, ?, ?, ?, ?);
    `, username, title, content, category, postTime)

	if err != nil {
		return 0, err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return postID, nil
}
