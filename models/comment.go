package models

import "time"

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
	Comment      Comment
	AuthorName   string
	LikeCount    int
	DislikeCount int
	Replies      []CommentView
	ReplyCount   int
	UserVote     *int
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
