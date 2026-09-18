package repository

import (
	"context"
	"database/sql"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	createUserPostReaction = ` INSERT INTO post_like (user_id, post_id, value) VALUES(?,?,?)`
	deleteUserPostReaction = ` DELETE FROM post_like WHERE (user_id = ? && post_id)`

	createUserCommentReaction = ` INSERT INTO comment_like (user_id, comment_id, value) VALUES(?,?,?)`
	deleteUserCommentReaction = ` DELETE FROM comment_like WHERE (user_id = ? && comment_id)`
)

// CREATE TABLE IF NOT EXISTS post_like (
//     user_id INTEGER NOT NULL,
//     post_id INTEGER NOT NULL,
//     value   INTEGER NOT NULL,
//     PRIMARY KEY (user_id, post_id),
//     FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
//     FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE
// );

// CREATE TABLE IF NOT EXISTS comment_like (
//     user_id    INTEGER NOT NULL,
//     comment_id INTEGER NOT NULL,
//     value      INTEGER NOT NULL,
//     PRIMARY KEY (user_id, comment_id),
//     FOREIGN KEY (user_id)    REFERENCES user(id)    ON DELETE CASCADE,
//     FOREIGN KEY (comment_id) REFERENCES comment(id) ON DELETE CASCADE
// );`

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) CreatePostReaction(cx context.Context, rec models.Reaction) error {

	_, postErr := r.db.ExecContext(cx, createUserPostReaction, rec.User_id, rec.Id, rec.Value)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return customerrors.ErrInternalError
	}
	return nil
}

func (r *SessionRepository) DeletePostReaction(cx context.Context, rec models.Reaction) error {
	_, deleteErr := r.db.ExecContext(cx, deleteUserPostReaction, rec.User_id, rec.Id)
	if deleteErr != nil {
		return customerrors.ErrInternalError
	}
	return nil
}

func (r *ReactionRepository) CreateCommentReaction(cx context.Context, rec models.Reaction) error {

	_, postErr := r.db.ExecContext(cx, createUserCommentReaction, rec.User_id, rec.Id, rec.Value)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return customerrors.ErrInternalError
	}
	return nil
}

func (r *SessionRepository) DeleteCommentReaction(cx context.Context, rec models.Reaction) error {
	_, deleteErr := r.db.ExecContext(cx, deleteUserPostReaction, rec.User_id, rec.Id)
	if deleteErr != nil {
		return customerrors.ErrInternalError
	}
	return nil
}
