package repository

import (
	"context"
	"database/sql"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

type commentRepositoryImpl struct {
	db *sql.DB
}

type CommentRepository interface {
	CreateComment(ctx context.Context, post models.CommentInput) (int, error)
	GetCommentsByPost(ctx context.Context, postID, int, userID *int) ([]*models.CommentView, error)
	GetRepliesByComment(ctx context.Context, commentID, int, userID *int) (*models.CommentView, error)
	UpdateComment(ctx context.Context, update models.CommentUpdate, authorID int) error
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	return &commentRepositoryImpl{
		db: db,
	}
}

func (r *commentRepositoryImpl) CreateComment(ctx context.Context, comment models.CommentInput) (int, error) {

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO comment (user_id, parent_post_id, parent_comment_id, content)
		VALUES (?, ?, ?, ?)
	`, comment.UserID, comment.ParentPostID, comment.ParentCommentID, comment.Content)
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	commentID, err := res.LastInsertId()
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	return int(commentID), nil
}

func (r *commentRepositoryImpl) GetCommentsByPost(ctx context.Context, postID, int, userID *int) ([]*models.CommentView, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.parent_post_id,
			c.parent_comment_id,
			c.content,
			c.created_at,
			c.updated_at,
			u.username,
			(
				SELECT COUNT(*)
				FROM comment_like cl
				WHERE cl.comment_id = c.id
				  AND cl.value = 1
			) AS like_count,
			(
				SELECT COUNT(*)
				FROM comment_like cl
				WHERE cl.comment_id = c.id
				  AND cl.value = -1
			) AS dislike_count,
			(
				SELECT COUNT(*)
				FROM comment reply
				WHERE reply.parent_comment_id = c.id
			) AS reply_count
		`
	args := []any{}

	if userID != nil {
		query += `, (
				SELECT cl.value
				FROM comment_like cl
				WHERE cl.comment_id = c.id
				AND cl.user_id = ?
			) AS user_vote
			`
		args = append(args, *userID)
	} else {
		query += `NULL AS user_vote`
	}

	// sort by oldest first by default, probably adding in toggle for oldest/newest/highest score first later
	query += `
		FROM comment c
		JOIN user u ON u.id = c.user_id
		WHERE c.parent_post_id = ?
		AND c.parent_comment_id IS NULL
		ORDER BY c.created_at ASC
		`
	args = append(args, postID)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	comments := []*models.CommentView{}

	for rows.Next() {
		var cv models.CommentView

		err := rows.Scan(
			&cv.Comment.ID,
			&cv.Comment.UserID,
			&cv.Comment.ParentPostID,
			&cv.Comment.ParentCommentID,
			&cv.Comment.Content,
			&cv.Comment.CreatedAt,
			&cv.Comment.UpdatedAt,
			&cv.AuthorName,
			&cv.LikeCount,
			&cv.DislikeCount,
			&cv.ReplyCount,
			&cv.UserVote,
		)
		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}

		cv.Replies = []models.CommentView{}
		comments = append(comments, &cv)
	}

	if err := rows.Err(); err != nil {
		return nil, customerrors.MapSQLError(err)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) GetRepliesByComment(ctx context.Context, commentID, int, userID *int) (*models.CommentView, error) {

}

func (r *commentRepositoryImpl) UpdateComment(ctx context.Context, update models.CommentUpdate, authorID int) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE comment
		SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
		`, update.Content, update.ID, authorID)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return customerrors.MapSQLError(err)
	}
	if rowCount == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}
