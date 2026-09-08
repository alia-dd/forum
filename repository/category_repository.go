package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

type categoryRepositoryImpl struct {
	db *sql.DB
}

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCatByPostIDs(ctx context.Context, postIDs []int) (map[int][]models.Category, error)
	CreateCategory(ctx context.Context, name string) (int, error)
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) CreateCategory(ctx context.Context, name string) (int, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO category
		(name) VALUES (?)`, name)
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, customerrors.MapSQLError(err)
	}

	return int(id), nil
}

func (r *categoryRepositoryImpl) GetCatByPostIDs(ctx context.Context, postIDs []int) (map[int][]models.Category, error) {
	res := make(map[int][]models.Category)
	if len(postIDs) == 0 {
		return res, nil
	}

	conds := make([]string, len(postIDs))
	args := make([]any, len(postIDs))
	for i, id := range postIDs {
		conds[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT pc.post_id, category.id, category.name
		FROM post_category pc
		JOIN category ON category.id = pc.category_id
		WHERE pc.post_id IN (%s)`, strings.Join(conds, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID int
		var cat models.Category
		if err := rows.Scan(&postID, &cat.ID, &cat.Name); err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		res[postID] = append(res[postID], cat)
	}
	return res, rows.Err()
}

func (r *categoryRepositoryImpl) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	ret := []models.Category{}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name
		FROM category
		ORDER BY name
	`)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		ret = append(ret, cat)
	}
	return ret, rows.Err()
}
