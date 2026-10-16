package repository

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const (
	RegisterUserQuery            = ` INSERT INTO user (username, name, email, bio, imagepath, password_hash) VALUES(?,?,?,?,?,?)`
	fetchUserInfoQuery           = ` SELECT id, username, name, email, ifnull(bio, ''), ifnull(imagepath, ''), password_hash, created_at, updated_at, role FROM user WHERE (username = ? or email = ?)`
	fetchUserInfoByIdQuery       = ` SELECT id, username, name, email, ifnull(bio, ''), ifnull(imagepath, ''), role FROM user WHERE id = ?`
	fetchUserInfoByUsernameQuery = ` SELECT id, username, name, ifnull(bio, ''), ifnull(imagepath, ''),  user.deleted_at IS NOT NULL FROM user WHERE username = ?`

	UpdateUserRoleByIdQuery = ` UPDATE user SET updated_at = CURRENT_TIMESTAMP, role = ? WHERE id = ?`

	UpdateUserByIdQuery = ` UPDATE user SET updated_at = CURRENT_TIMESTAMP, `
	DeleteUserQuery     = ` DELETE FROM user WHERE id = ?`

	fetchPasswordHashQuery         = ` SELECT password_hash FROM user WHERE id = ?`
	updatePasswordQuery            = ` UPDATE user SET  updated_at = CURRENT_TIMESTAMP, password_hash = ? WHERE id = ?`
	checkAvailableExcludeUserQuery = ` SELECT EXISTS(SELECT 1 FROM user WHERE (username = ? or email = ?) AND id != ?)`

	// admin specific
	fetchAllUsersQuery = ` SELECT  id, username, email, role, ifnull(imagepath, ''), created_at, updated_at  FROM user`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) RegisterUser(cx context.Context, u models.UserRegister) error {
	_, postErr := r.db.ExecContext(cx, RegisterUserQuery, u.Username, u.Name, u.Email, u.Bio, u.Image, u.Password)
	if postErr != nil {
		return customerrors.MapSQLError(postErr)
	}

	return nil
}

// check if user exists in the database
// this func is used for the login to authenticate the provided user credentials
func (r *UserRepository) AuthenticateUser(cx context.Context, col string) (models.UserInfo, error) {
	var user models.UserInfo
	var role int

	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoQuery, col, col).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Bio, &user.Image, &user.Password, &user.CreatedAt, &user.UpdatedAt, &role)
	if fetchErr != nil {
		return user, customerrors.MapSQLError(fetchErr)
	}

	switch role {
	case 700:
		user.Role = "admin"
	case 5:
		user.Role = "blocked"
	default:
		user.Role = "user"
	}

	return user, nil
}

func (r *UserRepository) FetchUserData(cx context.Context, userId int) (models.UserInfo, error) {
	var user models.UserInfo
	var role int

	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoByIdQuery, userId).Scan(&user.ID, &user.Username, &user.Name, &user.Email, &user.Bio, &user.Image, &role)
	if fetchErr != nil {
		return user, customerrors.MapSQLError(fetchErr)
	}

	switch role {
	case 700:
		user.Role = "admin"
	case 5:
		user.Role = "blocked"
	default:
		user.Role = "user"
	}

	return user, nil
}

// this repo func is used by the profile fetch for another user using their username
func (r *UserRepository) FetchUserDataByUserName(cx context.Context, username string) (models.PublicUserInfo, error) {
	var user models.PublicUserInfo
	fetchErr := r.db.QueryRowContext(cx, fetchUserInfoByUsernameQuery, username).Scan(&user.ID, &user.Username, &user.Name, &user.Bio, &user.Image, &user.DeletedAt)
	if fetchErr != nil {
		return user, customerrors.MapSQLError(fetchErr)
	}

	return user, nil
}

// this func implement a patch update on the user profile
func (r *UserRepository) UpdateUserData(cx context.Context, user models.UserUpdate, deleteAccount bool) error {
	var query string
	var extraQuery []string
	var args []any

	if deleteAccount {
		idString := strconv.Itoa(user.ID)
		extraQuery = append(extraQuery, "deleted_at = CURRENT_TIMESTAMP", "username = ?", "name = ?", "email = ?", "imagepath = ?", "bio = ?", "password_hash = ?")
		args = append(args, idString, "", idString, "", "", "")
	}

	if user.Username != nil {
		extraQuery = append(extraQuery, "username = ?")
		args = append(args, &user.Username)
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
	// fmt.Println(query)
	query = UpdateUserByIdQuery + strings.Join(extraQuery, ", ") + " WHERE id = ?"
	args = append(args, user.ID)
	_, err := r.db.ExecContext(cx, query, args...)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	return nil
}

func (r *UserRepository) UpdateUserRole(cx context.Context, user_id, role int) error {

	_, err := r.db.ExecContext(cx, UpdateUserRoleByIdQuery, role, user_id)
	if err != nil {
		return customerrors.MapSQLError(err)
	}

	return nil
}

func (r *UserRepository) DeleteUser(cx context.Context, userID int) error {
	res, deleteErr := r.db.ExecContext(cx, DeleteUserQuery, userID)
	if deleteErr != nil {
		return customerrors.MapSQLError(deleteErr)
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
		return "", customerrors.MapSQLError(fetchErr)
	}

	return hashPass, nil
}

func (r *UserRepository) UpdatePassword(cx context.Context, userID int, newPass string) error {
	_, updateErr := r.db.ExecContext(cx, updatePasswordQuery, newPass, userID)
	if updateErr != nil {
		return customerrors.MapSQLError(updateErr)
	}

	return nil
}

// this check is used by registration to see if the username/email is available
func (r *UserRepository) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user WHERE (username = ? or email = ?))`
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, query, col, col).Scan(&exist)

	if fetchErr != nil {
		return false, customerrors.MapSQLError(fetchErr)
	}

	return exist, nil
}

// this is used for the user profile edit to check if the profided username / email is available
// excluding the user incase they reuse the same email
func (r *UserRepository) CheckIfAvailableExcludeUser(cx context.Context, col string, userId int) (bool, error) {
	var exist bool
	fetchErr := r.db.QueryRowContext(cx, checkAvailableExcludeUserQuery, col, col, userId).Scan(&exist)
	if fetchErr != nil {
		return false, customerrors.MapSQLError(fetchErr)
	}

	return exist, nil
}

// admin specific query
func (r *UserRepository) FetchUsers(cx context.Context) ([]models.AdminUserInfo, error) {
	var users []models.AdminUserInfo

	rows, fetchErr := r.db.QueryContext(cx, fetchAllUsersQuery)
	if fetchErr != nil {
		return users, customerrors.MapSQLError(fetchErr)
	}

	for rows.Next() {
		var user models.AdminUserInfo
		err := rows.Scan(
			&user.Id, &user.Username, &user.Email,
			&user.Role, &user.Image, &user.CreatedAt, &user.UpdatedAt,
		)

		if err != nil {
			return nil, customerrors.MapSQLError(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, customerrors.MapSQLError(err)
	}

	return users, nil
}
