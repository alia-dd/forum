package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	RegisterUserQuery       = ` INSERT INTO user (username, name, email, password_hash) VALUES(?,?,?,?)`
	fetchUserInfoQurrey     = ` SELECT  id, username, name, email, password_hash, created_at, updated_at  FROM user WHERE (username = ? or email = ?)`
	fetchUserInfoByIdQurrey = ` SELECT id, username, name, email FROM user WHERE id = ?`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(cx context.Context, u models.UserRegister) error {
	_, postErr := r.db.ExecContext(cx, RegisterUserQuery, u.Username, u.Name, u.Email, u.Password)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to insert user  into table: %w", postErr)
	}

	return nil
}

// check if use exist in the database
func (r *UserRepository) AuthenticateUser(cx context.Context, col string) (models.UserInfo, error) {
	fmt.Println("here")
	var user models.UserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoQurrey, col, col).Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if fetchErr != nil {
		fmt.Println(fetchErr)
		if fetchErr == sql.ErrNoRows {
			return user, customerrors.ErrNotFound
		}
		return user, customerrors.ErrInternalError
	}
	return user, nil
}
func (r *UserRepository) FetchUserData(cx context.Context, userId int) (models.UserInfo, error) {

	var user models.UserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoByIdQurrey, userId).Scan(&user.Id, &user.Username, &user.Name, &user.Email)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return user, customerrors.ErrNotFound
		}
		return user, customerrors.ErrInternalError
	}
	return user, nil
}

func (r *UserRepository) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	qurrey := `SELECT EXISTS(SELECT 1 FROM user WHERE (username = ? or email = ?))`
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, qurrey, col, col).Scan(&exist)
	if fetchErr != nil {
		return false, fmt.Errorf("failed to fetch username from table: %w", fetchErr)
	}
	return exist, nil
}
