package repository

import (
	"context"
	"database/sql"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

type postRepositoryImpl struct {
	db *sql.DB
}

type PostRepository interface {
	CreatePost(ctx context.Context, post models.PostInput) (int, error)
	GetPost(ctx context.Context, filter models.PostFilter) ([]*models.PostView, error)
	GetPostByID(ctx context.Context, id int) (*models.PostView, error)
	UpdatePost(ctx context.Context, update models.PostUpdate, authorID int) error
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepositoryImpl{
		db: db,
	}
}

func (r *postRepositoryImpl) CreatePost(ctx context.Context, post models.PostInput) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		INSERT INTO post (user_id, title, content) VALUES (?, ?, ?)
	`, post.UserID, post.Title, post.Content)
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO post_category
		(post_id, category_id, is_main) VALUES (?, ?, 1)
		`, postID, post.MainCategoryID)
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	for _, catID := range post.CategoryIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO post_category (post_id, category_id) VALUES (?, ?)
		`, postID, catID)
		if err != nil {
			return 0, customerrors.MapSQLError(err)
		}
	}

	return int(postID), tx.Commit()
}

func (r *postRepositoryImpl) GetPost(ctx context.Context, filter models.PostFilter) ([]*models.PostView, error) {
	query := `
		SELECT
			post.id, post.user_id, post.title, post.content, post.created_at, post.updated_at,
			user.username,
			(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = 1)  AS like_count,
			(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = -1) AS dislike_count,
			(SELECT COUNT(*) FROM comment   WHERE comment.parent_post_id = post.id)                     AS comment_count
		FROM post
		JOIN user ON user.id = post.user_id
		`

	var conds []string
	var args []any

	if filter.AuthorID != nil {
		conds = append(conds, "post.user_id = ?")
		args = append(args, *filter.AuthorID)
	}

	if filter.CategoryID != nil {
		conds = append(conds, `
			EXISTS (
				SELECT 1 FROM post_category pc
				WHERE pc.post_id = post.id
				AND pc.category_id = ?
			)`)
		args = append(args, *filter.CategoryID)
	}

	if filter.LikedByID != nil {
		conds = append(conds, `
			EXISTS (
				SELECT 1 FROM post_like liked
				WHERE liked.post_id = post.id
				AND liked.user_id = ?
				AND liked.value = 1
			)`)
		args = append(args, *filter.LikedByID)
	}

	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	//order by matching main category, using creation time as tiebreak (main2 -> main1 -> side2 -> side1)
	if filter.CategoryID != nil {
		query += `
			ORDER BY (
				SELECT is_main FROM post_category
				WHERE post_id = post.id AND category_id = ?
			)   DESC, post.created_at DESC
		`
		args = append(args, *filter.CategoryID)
	} else {
		query += " ORDER BY post.created_at DESC"
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	posts := []*models.PostView{}

	for rows.Next() {
		var pv models.PostView

		err := rows.Scan(
			&pv.Post.ID,
			&pv.Post.UserID,
			&pv.Post.Title,
			&pv.Post.Content,
			&pv.Post.CreatedAt,
			&pv.Post.UpdatedAt,
			&pv.AuthorName,
			&pv.LikeCount,
			&pv.DislikeCount,
			&pv.CommentCount,
		)
		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}

		posts = append(posts, &pv)
	}

	if err := rows.Err(); err != nil {
		return nil, customerrors.MapSQLError(err)
	}

	return posts, nil
}

func (r *postRepositoryImpl) GetPostByID(ctx context.Context, id int) (*models.PostView, error) {
	query := `
		SELECT
			post.id, post.user_id, post.title, post.content, post.created_at, post.updated_at,
			user.username,
			(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = 1)  AS like_count,
			(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = -1) AS dislike_count,
			(SELECT COUNT(*) FROM comment   WHERE comment.parent_post_id = post.id)                     AS comment_count
		FROM post
		JOIN user ON user.id = post.user_id
		WHERE post.id = ?
		`

	var pv models.PostView
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pv.Post.ID, &pv.Post.UserID, &pv.Post.Title, &pv.Post.Content, &pv.Post.CreatedAt, &pv.Post.UpdatedAt,
		&pv.AuthorName, &pv.LikeCount, &pv.DislikeCount, &pv.CommentCount,
	)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	return &pv, nil
}

func (r *postRepositoryImpl) UpdatePost(ctx context.Context, update models.PostUpdate, authorID int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return customerrors.MapSQLError(err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `
		UPDATE post
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?
		`, update.Title, update.Content, update.ID, authorID)
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

	_, err = tx.ExecContext(ctx, `
		DELETE FROM post_category
		WHERE post_id = ?
		`, update.ID)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO post_category
		(post_id, category_id, is_main) VALUES (?, ?, 1)
		`, update.ID, update.MainCategoryID)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	for _, catID := range update.CategoryIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO post_category (post_id, category_id) VALUES (?, ?)
		`, update.ID, catID)
		if err != nil {
			return customerrors.MapSQLError(err)
		}
	}
	return tx.Commit()
}
