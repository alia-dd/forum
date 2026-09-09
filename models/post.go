package models

import (
	"time"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type PostView struct {
	Post         Post
	AuthorName   string
	Categories   []Category
	Comments     []CommentView // keep empty if loading a list of posts, load if viewing single post
	LikeCount    int
	DislikeCount int
	CommentCount int
	UserVote     *int
}

type PostInput struct {
	UserID         int
	Title          string
	Content        string
	MainCategoryID int
	CategoryIDs    []int
}

type PostUpdate struct {
	ID             int
	Title          string
	Content        string
	MainCategoryID int
	CategoryIDs    []int
}

type PostFilter struct {
	CategoryID *int
	AuthorID   *int
	LikedByID  *int
}
