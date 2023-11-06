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

	rows, err := db.Query("SELECT posts.id, posts.title, posts.content, posts.category, users.username FROM posts INNER JOIN users ON posts.user_id = users.id")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post

		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.Username); err != nil {

			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func InsertComment(db *sql.DB, postID, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (user_id, post_id, content, creation_date) VALUES (?, ?, ?, ?)",
		userID, postID, content, time.Now())
	return err
}
func GetCommentsForPost(db *sql.DB, postID int) ([]Comment, error) {
	var comments []Comment
	rows, err := db.Query("SELECT comments.id, comments.post_id, comments.user_id, comments.content, comments.created_at, users.username FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, comment.Content, &comment.CreationDate, &comment.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func InsertPostLike(db *sql.DB, userID, postID int) error {
	var likeStatus bool

	err := db.QueryRow("SELECT like_status FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&likeStatus)
	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO post_likes (user_id, post_id, like_status, like) VALUES (?, ?, ?, ?)", userID, postID, true, 1)
		return err
	} else if err != nil {
		return err
	} else {
		if likeStatus {
			_, err := db.Exec("DELETE FROM post_likes WHERE user_id = ? and post_id = ?", userID, postID)
			return err
		} else {
			_, err := db.Exec("UPDATE post_likes SET like_status = ?, like = like + 1 WHERE user_id = ? AND post_id = ?", true, userID, postID)
			return err
		}
	}
}

// func InsertDislikePost(db *sql.DB, userID, postID int) error {
// 	var likeStatus bool
// 	err := db.Query("SELECT like_status FROM post_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan()
// 	return err
// }
