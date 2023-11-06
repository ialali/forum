package database

type User struct {
	ID               int
	Username         string
	Email            string
	Password         string
	RegistrationDate string
}

type Post struct {
	ID           int
	UserID       int
	Title        string
	Content      string
	Category     string
	Comments     []Comment
	CreationDate string
	Username     string
}
type Comment struct {
	ID           int
	UserID       int
	PostID       int
	Username     string
	Content      string
	CreationDate string
}
type Category struct {
	ID   int
	Name string
}
type PostLike struct {
	ID     int
	UserID int
	PostID int
	Like   int // 1 for like, 0 for dislike
}

type CommentLike struct {
	ID        int
	UserID    int
	CommentID int
	Like      int // 1 for like, 0 for dislike
}
