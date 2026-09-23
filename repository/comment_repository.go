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
	GetCommentsByPost(ctx context.Context, postID int, userID *int) ([]*models.CommentView, error)
	GetCommentByID(ctx context.Context, comID int, curUserID *int) (*models.CommentView, error)
	// // GetRepliesByComment(ctx context.Context, commentID, int, userID *int) (*models.CommentView, error)
	UpdateCommentContent(ctx context.Context, commentID, userID int, content string) error
	CountByPost(ctx context.Context, postID int) (int, error)
	DeleteComment(ctx context.Context, id, authorID int) error
	GetRepliesByCommentID(ctx context.Context, parentCommentID int, userID *int) ([]*models.CommentView, error)
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

func (r *commentRepositoryImpl) GetCommentByID(ctx context.Context, comID int, curUserID *int) (*models.CommentView, error) {
	query := `
        SELECT c.id, c.parent_post_id, c.parent_comment_id, c.user_id, c.content, c.created_at, u.username
        FROM comment c
        JOIN user u ON c.user_id = u.id
        WHERE c.id = ?
    `

	comment := &models.CommentView{}
	err := r.db.QueryRowContext(ctx, query, comID).Scan(
		&comment.ID,
		&comment.ParentPostID,
		&comment.ParentCommentID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
		&comment.AuthorName,
	)
	if err != nil {
		return nil, err
	}

	comment.IsOwner = false
	comment.IsLogged = false

	if curUserID != nil {
		comment.IsLogged = true
		if comment.UserID == *curUserID {
			comment.IsOwner = true
		}
	}

	comment.LikeCount = 0
	comment.DislikeCount = 0
	comment.ReplyCount = 0
	comment.Replies = nil
	comment.UserVote = nil

	return comment, nil
}

func (r *commentRepositoryImpl) GetCommentsByPost(ctx context.Context, postID int, curUserID *int) ([]*models.CommentView, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.parent_post_id,
			c.parent_comment_id,
			c.content,
			c.created_at,
			c.updated_at,
			c.deleted_at,
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

	if curUserID != nil {
		query += `, (
				SELECT cl.value
				FROM comment_like cl
				WHERE cl.comment_id = c.id
				AND cl.user_id = ?
			) AS user_vote
			`
		args = append(args, *curUserID)
	} else {
		query += `, NULL AS user_vote`
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
		var deletedAt sql.NullTime

		err := rows.Scan(
			&cv.Comment.ID,
			&cv.Comment.UserID,
			&cv.Comment.ParentPostID,
			&cv.Comment.ParentCommentID,
			&cv.Comment.Content,
			&cv.Comment.CreatedAt,
			&cv.Comment.UpdatedAt,
			&deletedAt,
			&cv.AuthorName,
			&cv.LikeCount,
			&cv.DislikeCount,
			&cv.ReplyCount,
			&cv.UserVote,
		)
		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}

		cv.Deleted = deletedAt.Valid

		cv.IsOwner = false
		cv.IsLogged = false

		if curUserID != nil {
			cv.IsLogged = true
			if cv.UserID == *curUserID {
				cv.IsOwner = true
			}
		}

		cv.Replies = []models.CommentView{}
		comments = append(comments, &cv)
	}

	if err := rows.Err(); err != nil {
		return nil, customerrors.MapSQLError(err)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) UpdateCommentContent(ctx context.Context, commentID, userID int, content string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE comment
		SET content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
		`, content, commentID, userID)
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

func (r *commentRepositoryImpl) CountByPost(ctx context.Context, postID int) (int, error) {
	query := `SELECT COUNT(*) FROM comment WHERE parent_post_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *commentRepositoryImpl) DeleteComment(ctx context.Context, id, authorID int) error {
	var hasReplies bool
	var err error
	err = r.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM comment
			WHERE parent_comment_id = ?)
		`, id).Scan(&hasReplies)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	var res sql.Result
	if !hasReplies {
		res, err = r.db.ExecContext(ctx, `
			DELETE FROM comment
			WHERE id = ? AND user_id = ?
			`, id, authorID)
	} else {
		res, err = r.db.ExecContext(ctx, `
			UPDATE comment
			SET content = '', deleted_at = CURRENT_TIMESTAMP
			WHERE id = ? AND user_id = ?
			`, id, authorID)
	}
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

func (r *commentRepositoryImpl) GetRepliesByCommentID(ctx context.Context, parentCommentID int, curUserID *int) ([]*models.CommentView, error) {
    query := `
        SELECT
            c.id,
            c.user_id,
            c.parent_post_id,
            c.parent_comment_id,
            c.content,
            c.created_at,
            c.updated_at,
            c.deleted_at,
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

    if curUserID != nil {
        query += `, (
                SELECT cl.value
                FROM comment_like cl
                WHERE cl.comment_id = c.id
                AND cl.user_id = ?
            ) AS user_vote
            `
        args = append(args, *curUserID)
    } else {
        query += `, NULL AS user_vote`
    }

    // Fetch comments belonging to this specific parent comment
    query += `
        FROM comment c
        JOIN user u ON u.id = c.user_id
        WHERE c.parent_comment_id = ?
        ORDER BY c.created_at ASC
        `
    args = append(args, parentCommentID)

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, customerrors.MapSQLError(err)
    }
    defer rows.Close()

    comments := []*models.CommentView{}

    for rows.Next() {
        var cv models.CommentView
        var deletedAt sql.NullTime

        err := rows.Scan(
            &cv.Comment.ID,
            &cv.Comment.UserID,
            &cv.Comment.ParentPostID,
            &cv.Comment.ParentCommentID,
            &cv.Comment.Content,
            &cv.Comment.CreatedAt,
            &cv.Comment.UpdatedAt,
            &deletedAt,
            &cv.AuthorName,
            &cv.LikeCount,
            &cv.DislikeCount,
            &cv.ReplyCount,
            &cv.UserVote,
        )
        if err != nil {
            return nil, customerrors.MapSQLError(err)
        }

        cv.Deleted = deletedAt.Valid

        cv.IsOwner = false
		cv.IsLogged = false

		if curUserID != nil {
			cv.IsLogged = true
			if cv.UserID == *curUserID {
				cv.IsOwner = true
			}
		}

        cv.Replies = []models.CommentView{}
        comments = append(comments, &cv)
    }

    if err := rows.Err(); err != nil {
        return nil, customerrors.MapSQLError(err)
    }

    return comments, nil
}