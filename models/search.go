package models

type UserResult struct {
	ID       int
	Username string
	Name     string
}

type PostResult struct {
	ID         int
	Title      string
	Content    string
	AuthorName string
}

type CommentResult struct {
	ID           int
	Content      string
	ParentPostID int
	AuthorName   string
}
