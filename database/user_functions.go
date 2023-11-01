package database

import (
	"database/sql"
	"log"

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
func GetIDBYusername(db *sql.DB, username string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CreatePost inserts a new post into the database and returns the post ID.
func InsertPost(db *sql.DB, category, title, content string, userID int) error {
	// You should also add additional error handling and validation as needed.

	// Prepare the SQL statement to insert a new post.
	stmt, err := db.Prepare("INSERT INTO posts (title, content, user_id, category) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Get the current timestamp.
	createdAt := time.Now()

	// Execute the SQL statement to insert the new post.
	_, err = stmt.Exec(title, content, userID, createdAt)
	if err != nil {
		return err
	}

	return nil
}

func GetPosts(db *sql.DB) ([]Post, error) {
	var posts []Post

	// Query to retrieve posts
	rows, err := db.Query("SELECT id, title, content, user_id FROM posts")
	if err != nil {
		log.Printf("Error querying posts: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error reading rows: %v\n", err)
		return nil, err
	}

	log.Printf("Fetched %d posts from the database.\n", len(posts))
	return posts, nil
}
