package models

import (
	"time"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	Image     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type PostView struct {
	Post
	AuthorName    string
	AuthorImage   string
	AuthorDeleted bool
	Categories    []Category
	Comments      []CommentView // keep empty if loading a list of posts, load if viewing single post
	LikeCount     int
	DislikeCount  int
	UserReaction  int // user post Reaction value
	TargetType    string
	TargetId      int
	CommentCount  int
	Deleted       bool
}

type PostInput struct {
	UserID         int
	Title          string
	Content        string
	Image          string
	MainCategoryID int
	CategoryIDs    []int
}

type PostUpdate struct {
	ID             int
	Title          string
	Content        string
	Image          string
	MainCategoryID int
	CategoryIDs    []int
}

type PostFilter struct {
	CategoryID *int
	AuthorID   *int
	LikedByID  *int
	Limit      int
	Offset     int
}

func (p Post) CreatedAgo() string {
	return timeAgo(p.CreatedAt)
}

func (p Post) UpdatedAgo() string {
	if p.UpdatedAt == nil {
		return ""
	}
	return timeAgo(*p.UpdatedAt)
}
