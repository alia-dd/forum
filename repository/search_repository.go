package repository

import (
	"context"
	"database/sql"
	"fmt"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) SearchUsers(ctx context.Context, search string) ([]models.UserResult, error) {
	// relevance/newest sort does nothing here, minor UI issue?
	searchArg := "%" + search + "%"
	prefixArg := search + "%"

	rows, err := r.db.QueryContext(ctx, `
			SELECT id, username
			FROM user
			WHERE username LIKE ?
			ORDER BY
				CASE
					WHEN username LIKE ? THEN 0 
					ELSE 1
				END,
				username
			`, searchArg, prefixArg) // order results so that matches that begin with the searchArg are shown before those that contain it 'in' them
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	result := []models.UserResult{}
	for rows.Next() {
		var ur models.UserResult
		err := rows.Scan(&ur.ID, &ur.Username)
		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		result = append(result, ur)
	}
	return result, rows.Err()
}

func (r *SearchRepository) SearchPosts(ctx context.Context, term string, sort string) ([]models.PostResult, error) {
	var orderBy string

	switch sort {
	case "time":
		orderBy = "post.created_at DESC"
	default:
		orderBy = "bm25(post_fts, 10.0, 1.0), post.created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT post.id, post.title, user.username
		FROM post_fts
		JOIN post ON post.id = post_fts.rowid
		JOIN user ON user.id = post.user_id
		WHERE post_fts MATCH ?
		ORDER BY %s
		`, orderBy)

	// not checking for if post is deleted, since title should still be searchable even then
	rows, err := r.db.QueryContext(ctx, query, term)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	result := []models.PostResult{}
	for rows.Next() {
		var pr models.PostResult
		if err := rows.Scan(&pr.ID, &pr.Title, &pr.AuthorName); err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		result = append(result, pr)
	}
	return result, rows.Err()
}

func (r *SearchRepository) SearchComments(ctx context.Context, term string, sort string) ([]models.CommentResult, error) {
	var orderBy string

	switch sort {
	case "time":
		orderBy = "comment.created_at DESC"
	default:
		orderBy = "bm25(comment_fts), comment.created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT comment.id, comment.content, comment.parent_post_id, user.username
		FROM comment_fts
		JOIN comment ON comment.id = comment_fts.rowid
		JOIN user ON user.id = comment.user_id
		WHERE comment_fts MATCH ?
		AND comment.deleted_at IS NULL
		ORDER BY %s
		`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, term)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	result := []models.CommentResult{}
	for rows.Next() {
		var c models.CommentResult
		if err := rows.Scan(&c.ID, &c.Content, &c.ParentPostID, &c.AuthorName); err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *SearchRepository) SearchAdminUsers(ctx context.Context, search string) ([]models.AdminUserInfo, error) {
	searchArg := "%" + search + "%"
	prefixArg := search + "%"

	rows, err := r.db.QueryContext(ctx, `
			SELECT id, username, email, role, COALESCE(imagepath, '')
			FROM user
			WHERE username LIKE ?
			ORDER BY
				CASE
					WHEN username LIKE ? THEN 0
					ELSE 1
				END,
				username
			`, searchArg, prefixArg)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	result := []models.AdminUserInfo{}
	for rows.Next() {
		var ur models.AdminUserInfo
		err := rows.Scan(&ur.Id, &ur.Username, &ur.Email, &ur.Role, &ur.Image)
		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		result = append(result, ur)
	}
	return result, rows.Err()
}

func (r *SearchRepository) SearchAdminPosts(ctx context.Context, term string) ([]*models.PostView, error) {
	orderBy := "bm25(post_fts, 10.0, 1.0), post.created_at DESC"

	query := fmt.Sprintf(`
		SELECT post.id, post.user_id, post.title, post.content, post.created_at, post.updated_at, user.username,
		(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = 1)  AS like_count,
			(SELECT COUNT(*) FROM post_like WHERE post_like.post_id = post.id AND post_like.value = -1) AS dislike_count,
			(SELECT COUNT(*) FROM comment   WHERE comment.parent_post_id = post.id)                     AS comment_count
		FROM post_fts
		JOIN post ON post.id = post_fts.rowid
		JOIN user ON user.id = post.user_id
		WHERE post_fts MATCH ?
		ORDER BY %s
		`, orderBy)

	rows, err := r.db.QueryContext(ctx, query, term)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	result := []*models.PostView{}
	for rows.Next() {
		pv := &models.PostView{}
		if err := rows.Scan(
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
		); err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		result = append(result, pv)
	}
	return result, rows.Err()
}
