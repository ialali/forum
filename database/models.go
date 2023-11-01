package database

type User struct {
	ID               int
	Username         string
	Email            string
	Password         string
	RegistrationDate string
}

type Post struct {
	ID       int
	UserID   int
	Title    string
	Content  string
	Category string

	CreationDate string
	Username     string
}
type Comment struct {
	ID           int
	UserID       int
	PostID       int
	Content      string
	CreationDate string
}
type Category struct {
	ID   int
	Name string
}
type Like struct {
	ID       int
	UserID   int
	ItemID   int
	ItemType string
}
