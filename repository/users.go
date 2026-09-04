package repository

import (
	"context"
	"database/sql"
	"fmt"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreatUser(cx context.Context, u models.UserRegister) error {
	qurrey := ` INSERT INTO users (username,email,password_hash) VALUE(?,?,?)`
	_, postErr := r.db.ExecContext(cx, qurrey, u.Username, u.Email, u.Password)
	if postErr != nil {
		return fmt.Errorf("failed to insert suer into table: %w", postErr)
	}
	return nil
}

func (r *UserRepository) GetUserById(cx context.Context, id int) error {
	qurrey := ` SELECT * FROM users`
	_, fetchErr := r.db.QueryContext(cx, qurrey, id)
	if fetchErr != nil {
		return fmt.Errorf("failed to fetch suer info from table: %w", fetchErr)
	}
	return nil
}

func (r *UserRepository) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	qurrey := `SELECT EXISTS(SELECT 1 FROM users WHERE (username = ? or email = ?))`
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, qurrey, col, col).Scan(&exist)
	if fetchErr != nil {
		return false, fmt.Errorf("failed to fetch username from table: %w", fetchErr)
	}
	return !exist, nil
}
