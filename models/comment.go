package models

import (
	"time"
)

type Comment struct {
	ID              int
	UserID          int
	ParentPostID    int
	ParentCommentID *int
	Content         string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

type CommentView struct {
	Comment
	AuthorName   string
	LikeCount    int
	DislikeCount int
	Replies      []CommentView
	ReplyCount   int
	IsOwner      bool
	IsLogged     bool
	Deleted      bool
	UserReaction int // user post Reaction value
	TargetType   string
	TargetId     int
}

type CommentInput struct {
	UserID          int
	ParentPostID    int
	ParentCommentID *int
	Content         string
}

type CommentUpdate struct {
	ID      int
	Content string
}
