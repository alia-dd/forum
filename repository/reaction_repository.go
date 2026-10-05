package repository

import (
	"context"
	"database/sql"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	getUserPostReaction    = ` SELECT value FROM post_like WHERE user_id = ? AND post_id = ?`
	createUserPostReaction = ` INSERT INTO post_like (user_id, post_id, value) VALUES(?,?,?)`
	deleteUserPostReaction = ` DELETE FROM post_like WHERE user_id = ? AND post_id = ?`
	updateUserPostReaction = ` UPDATE post_like SET value = ? WHERE user_id = ? AND post_id = ?`

	getUserCommentReaction    = ` SELECT value FROM comment_like WHERE user_id = ? AND comment_id = ?`
	createUserCommentReaction = ` INSERT INTO comment_like (user_id, comment_id, value) VALUES(?,?,?)`
	deleteUserCommentReaction = ` DELETE FROM comment_like WHERE user_id = ? AND comment_id = ?`
	updateUserCommentReaction = ` UPDATE comment_like SET value = ? WHERE user_id = ? AND comment_id = ?`
)

type ReactionRepository struct {
	db *sql.DB
}

func NewReactionRepository(db *sql.DB) *ReactionRepository {
	return &ReactionRepository{db: db}
}

func (r *ReactionRepository) GetPostReaction(cx context.Context, rec models.Reaction) (*int, error) {
	var value int
	fetchErr := r.db.QueryRowContext(cx, getUserPostReaction, rec.User_id, rec.ID).Scan(&value)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, customerrors.ErrInternalError
	}
	return &value, nil
}
func (r *ReactionRepository) CreatePostReaction(cx context.Context, rec models.Reaction) error {
	_, postErr := r.db.ExecContext(cx, createUserPostReaction, rec.User_id, rec.ID, rec.Value)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return customerrors.ErrInternalError
	}
	return nil
}
func (r *ReactionRepository) UpdatePostReaction(cx context.Context, rec models.Reaction) error {
	_, err := r.db.ExecContext(cx, updateUserPostReaction, rec.Value, rec.User_id, rec.ID)

	if err != nil {
		return customerrors.ErrInternalError
	}
	return nil
}
func (r *ReactionRepository) DeletePostReaction(cx context.Context, rec models.Reaction) error {
	_, deleteErr := r.db.ExecContext(cx, deleteUserPostReaction, rec.User_id, rec.ID)
	if deleteErr != nil {
		return customerrors.ErrInternalError
	}
	return nil
}

// comment reaction repo
func (r *ReactionRepository) GetCommentReaction(cx context.Context, rec models.Reaction) (*int, error) {
	var value int
	fetchErr := r.db.QueryRowContext(cx, getUserCommentReaction, rec.User_id, rec.ID).Scan(&value)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return nil, nil
		}
		return nil, customerrors.ErrInternalError
	}
	return &value, nil
}

func (r *ReactionRepository) CreateCommentReaction(cx context.Context, rec models.Reaction) error {
	_, postErr := r.db.ExecContext(cx, createUserCommentReaction, rec.User_id, rec.ID, rec.Value)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return customerrors.ErrInternalError
	}
	return nil
}

func (r *ReactionRepository) UpdateCommentReaction(cx context.Context, rec models.Reaction) error {
	_, err := r.db.ExecContext(cx, updateUserCommentReaction, rec.Value, rec.User_id, rec.ID)
	if err != nil {
		return customerrors.ErrInternalError
	}

	return nil
}

func (r *ReactionRepository) DeleteCommentReaction(cx context.Context, rec models.Reaction) error {
	_, deleteErr := r.db.ExecContext(cx, deleteUserCommentReaction, rec.User_id, rec.ID)
	if deleteErr != nil {
		return customerrors.ErrInternalError
	}
	return nil
}
