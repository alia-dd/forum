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
	LikeCount    int
	DislikeCount int
	CommentCount int
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
