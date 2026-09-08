package models
<<<<<<< HEAD
=======

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

type PostUpdate struct {
	ID          int
	Title       string
	Content     string
	CategoryIDs []int
}

type PostView struct {
	Post         Post
	AuthorName   string
	Categories   []Category
	LikeCount    int
	DislikeCount int
	CommentCount int
}

type PostInput struct {
	UserID      int
	Title       string
	Content     string
	CategoryIDs []int
}

type PostFilter struct {
	CategoryID *int
	AuthorID   *int
	LikedByID  *int
}
>>>>>>> main
