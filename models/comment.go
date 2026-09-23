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
	UserVote     *int
	IsOwner      bool
	Deleted      bool
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
