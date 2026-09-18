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
	RegisterUserQuery              = ` INSERT INTO user (username, name, email, imagepath, bio, password_hash) VALUES(?,?,?,?,?,?)`
	fetchUserInfoQurrey            = ` SELECT  id, username, name, email, bio, imagepath, password_hash, created_at, updated_at  FROM user WHERE (username = ? or email = ?)`
	fetchUserInfoByIdQurrey        = ` SELECT id, username, name, email, bio, imagepath FROM user WHERE id = ?`
	fetchUserInfoByUsernameQurey   = ` SELECT id, username, name, ifnull(bio, ''), ifnull(imagepath, '')  FROM user WHERE username = ?`
	UpdateUserByIdQuery            = ` UPDATE user SET updated_at = CURRENT_TIMESTAMP, `
	DeleteUserQuery                = ` DELETE user WHERE id = ?`
	fetchPasswordHashQuery         = ` SELECT password_hash FROM user WHERE id = ?`
	updatePasswordQuery            = ` UPDATE user SET  updated_at = CURRENT_TIMESTAMP, password_hash = ? WHERE id = ?`
	checkAvailableExcludeUserQuery = ` SELECT EXISTS(SELECT 1 FROM user WHERE (username = ? or email = ?) AND id != ?)`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(cx context.Context, u models.UserRegister) error {
	_, postErr := r.db.ExecContext(cx, RegisterUserQuery, u.Username, u.Name, u.Email, u.Image, u.Bio, u.Password)
	if postErr != nil {
		if strings.Contains(postErr.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrDuplicateEntry
		}
		return fmt.Errorf("failed to insert user  into table: %w", postErr)
	}

	return nil
}

// check if use exist in the database
// this func is used for the login to authentica the profided use credentials
func (r *UserRepository) AuthenticateUser(cx context.Context, col string) (models.UserInfo, error) {
	var user models.UserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoQurrey, col, col).Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.Bio, &user.Image, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return user, customerrors.ErrInvalidLogin
		}
		return user, customerrors.ErrInternalError
	}
	return user, nil
}

func (r *UserRepository) FetchUserData(cx context.Context, userId int) (models.UserInfo, error) {
	var user models.UserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoByIdQurrey, userId).Scan(&user.Id, &user.Username, &user.Name, &user.Email, &user.Bio, &user.Image)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return user, customerrors.ErrNotFound
		}
		return user, customerrors.ErrInternalError
	}
	return user, nil
}

// this repo func is used by the profile fetch for another user using their username
func (r *UserRepository) FetchUserDataByUserName(cx context.Context, username string) (models.PublicUserInfo, error) {
	var user models.PublicUserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoByUsernameQurey, username).Scan(&user.Id, &user.Username, &user.Name, &user.Bio, &user.Image)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return user, customerrors.ErrNotFound
		}
		return user, customerrors.ErrInternalError
	}
	return user, nil
}

// this func implement a patch update on the user profile
func (r *UserRepository) UpdateUserData(cx context.Context, user models.UserUpdate) error {

	var extraQuery []string
	var args []any

	if user.Username != nil {
		extraQuery = append(extraQuery, "username = ?")
		args = append(args, *user.Username)
	}

	if user.Name != nil {
		extraQuery = append(extraQuery, "name = ?")
		args = append(args, *user.Name)
	}

	if user.Email != nil {
		extraQuery = append(extraQuery, "email = ?")
		args = append(args, *user.Email)
	}
	if user.Bio != nil {
		extraQuery = append(extraQuery, "bio = ?")
		args = append(args, *user.Bio)
	}
	if user.Image != nil {
		extraQuery = append(extraQuery, "imagepath = ?")
		args = append(args, *user.Image)
	}

	if len(extraQuery) == 0 {
		return nil
	}

	query := UpdateUserByIdQuery + strings.Join(extraQuery, ", ") + " WHERE id = ?"

	args = append(args, user.Id)
	_, err := r.db.ExecContext(cx, query, args...)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return customerrors.ErrBadRequest
		}

		return customerrors.ErrInternalError
	}

	return nil
}

func (r *UserRepository) DeleteUser(cx context.Context, userID int) error {
	res, deleteErr := r.db.ExecContext(cx, DeleteUserQuery, userID)
	if deleteErr != nil {
		return customerrors.ErrInternalError
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return customerrors.ErrNotFound
	}
	return nil
}

// password reset repo functions
func (r *UserRepository) GetPasswordHash(cx context.Context, userID int) (string, error) {
	var hashPass string
	fetchErr := r.db.QueryRowContext(cx, fetchPasswordHashQuery, userID).Scan(&hashPass)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			return "", customerrors.ErrNotFound
		}
		return "", customerrors.ErrInternalError
	}

	return hashPass, nil
}

func (r *UserRepository) UpdatePassword(cx context.Context, userID int, newPass string) error {
	_, fetchErr := r.db.ExecContext(cx, updatePasswordQuery, newPass, userID)
	if fetchErr != nil {
		return customerrors.ErrInternalError
	}
	return nil
}

// this check if used by registrasion to see if the username/email are availabe
func (r *UserRepository) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	qurrey := `SELECT EXISTS(SELECT 1 FROM user WHERE (username = ? or email = ?))`
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, qurrey, col, col).Scan(&exist)
	if fetchErr != nil {
		return false, fmt.Errorf("failed to fetch username from table: %w", fetchErr)
	}
	return exist, nil
}

// this is used for the use profile edit to check if the profided username / email is availabe
// excluding the user incase they reuse the same email
func (r *UserRepository) CheckIfAvailableExcludeUser(cx context.Context, col string, userId int) (bool, error) {
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, checkAvailableExcludeUserQuery, col, col, userId).Scan(&exist)
	if fetchErr != nil {
		return false, fmt.Errorf("failed to fetch username from table: %w", fetchErr)
	}
	return exist, nil
}
