package repository

import (
	"context"
	"database/sql"

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

func (r *SearchRepository) SearchPosts(ctx context.Context, term string) ([]models.PostResult, error) {
	searchArg := "%" + term + "%"

	rows, err := r.db.QueryContext(ctx, `
		SELECT post.id, post.title, user.username
		FROM post JOIN user ON user.id = post.user_id
		WHERE post.title LIKE ? OR post.content LIKE ?
		ORDER BY post.created_at DESC
		`, searchArg, searchArg)
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

func (r *SearchRepository) SearchComments(ctx context.Context, term string) ([]models.CommentResult, error) {
	searchArg := "%" + term + "%"

	rows, err := r.db.QueryContext(ctx, `
		SELECT comment.id, comment.content, comment.parent_post_id, user.username
		FROM comment JOIN user ON user.id = comment.user_id
		WHERE comment.content LIKE ? AND comment.deleted_at IS NULL
		ORDER BY comment.created_at DESC
		`, searchArg)
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
