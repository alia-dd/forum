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
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCatByPostIDs(ctx context.Context, postIDs []int) (map[int][]models.Category, error)
	CreateCategory(ctx context.Context, name string) (int, error)
	UpdateCategory(ctx context.Context, id int, name string) error
	DeleteCategory(ctx context.Context, id int) error
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
	}
}

func (r *categoryRepositoryImpl) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	var cat models.Category
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name
		FROM category
		WHERE id = ?
		`, id).Scan(&cat.ID, &cat.Name)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	return &cat, nil
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
		SELECT pc.post_id, category.id, category.name, pc.is_main
		FROM post_category pc
		JOIN category ON category.id = pc.category_id
		WHERE pc.post_id IN (%s)
		ORDER BY pc.post_id, pc.is_main DESC, category.name`, strings.Join(conds, ", "))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, customerrors.MapSQLError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var postID int
		var cat models.Category
		if err := rows.Scan(&postID, &cat.ID, &cat.Name, &cat.IsMain); err != nil {
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

func (r *categoryRepositoryImpl) UpdateCategory(ctx context.Context, id int, name string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE category
		SET name = ? WHERE id = ?
		`, name, id)
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

func (r *categoryRepositoryImpl) DeleteCategory(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return customerrors.MapSQLError(err)
	}
	defer tx.Rollback()

	var used bool
	err = tx.QueryRowContext(ctx, `
		SELECT EXISTS
		(SELECT 1 FROM post_category
		WHERE category_id = ?)
		`, id).Scan(&used)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	if used { // no force deletion on purpose
		return customerrors.ErrInUse
	}

	res, err := tx.ExecContext(ctx, `
		DELETE FROM category
		WHERE id = ?
		`, id)
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
	return tx.Commit()
}
