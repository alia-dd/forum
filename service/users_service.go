package service

import (
	"context"
	"fmt"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type UserService struct {
	repo        *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewUserService(repo *repository.UserRepository, sessionRepo *repository.SessionRepository) *UserService {
	return &UserService{repo: repo, sessionRepo: sessionRepo}
}

// register new use
func (s *UserService) CreateUserService(cx context.Context, u models.UserRegister) error {

	if validationsErr := u.Isvalid(); validationsErr != nil {
		return validationsErr
	}
	hashedPass, hashErr := utils.GenerateHashPassword(u.Password)
	if hashErr != nil {
		return customerrors.ErrInternalError
	}

	u.Password = string(hashedPass)
	return s.repo.RegisterUser(cx, u)
}

func (s *UserService) GetUserService(cx context.Context, username string) (models.PublicUserInfo, error) {
	return s.repo.FetchUserDataByUserName(cx, username)
}

// fetch all user
func (s *UserService) GetAllUsersService(cx context.Context) ([]models.AdminUserInfo, error) {
	return s.repo.FetchUsers(cx)
}

// patch user
func (s *UserService) UpdateUserService(cx context.Context, u models.UserUpdate) error {

	if validationsErr := u.Isvalid(); validationsErr != nil {
		return validationsErr
	}
	if u.Username != nil {
		notAvailable, availableErr := s.repo.CheckIfAvailableExcludeUser(cx, *u.Username, u.ID)
		if availableErr != nil {
			return availableErr
		}
		if notAvailable {
			return customerrors.ErrDuplicateEntry
		}
	}

	if u.Email != nil {
		notAvailable, availableErr := s.repo.CheckIfAvailableExcludeUser(cx, *u.Email, u.ID)
		if availableErr != nil {
			return availableErr
		}
		if notAvailable {
			return customerrors.ErrDuplicateEntry
		}
	}
	return s.repo.UpdateUserData(cx, u)
}

func (s *UserService) UpdateUserRoleService(cx context.Context, action string, user_id int) error {
	var err error

	switch action {
	case "role":
		role := 0
		userData, fetchErr := s.repo.FetchUserData(cx, user_id)
		if fetchErr != nil {
			return fetchErr
		}
		if userData.Role == "admin" || userData.Role == "blocked" {
			role = 0
		} else {
			role = 700
		}
		fmt.Println(role)
		err = s.repo.UpdateUserRole(cx, user_id, role)
	case "block":
		role := 5
		userData, fetchErr := s.repo.FetchUserData(cx, user_id)
		if fetchErr != nil {
			return fetchErr
		}
		if userData.Role == "blocked" {
			role = 0
		}
		fmt.Println(role)
		err = s.repo.UpdateUserRole(cx, user_id, role)
	case "delete":
		err = s.repo.DeleteUser(cx, user_id)
		fmt.Println(err)
	default:
		return fmt.Errorf("unknown user action: %s", action)
	}
	return err
}

func (s *UserService) ChangePasswordService(cx context.Context, userID int, currentPass, newPass string) error {
	currentPassHash, fetchErr := s.repo.GetPasswordHash(cx, userID)
	if fetchErr != nil {
		return fetchErr
	}
	match, err := utils.VerifyPassword(currentPassHash, currentPass)
	if err != nil {
		return customerrors.ErrInternalError
	}
	if !match {
		return customerrors.ErrIncorrectPassword //401
	}

	newPassHash, hashErr := utils.GenerateHashPassword(newPass)
	if hashErr != nil {
		return customerrors.ErrInternalError
	}

	if updateErr := s.repo.UpdatePassword(cx, userID, string(newPassHash)); updateErr != nil {
		fmt.Println(">>>", updateErr)
		return updateErr
	}
	return s.sessionRepo.DeleteSessionsByUserId(cx, userID)
}

// authernticate use with their email/usename and password
func (s *UserService) AuthenticateUserService(cx context.Context, u models.UserLogin) (*models.Session, error) {

	user, fetchErr := s.repo.AuthenticateUser(cx, u.Username)
	if fetchErr != nil {
		return nil, fetchErr
	}
	if user.Role == "blocked" {
		return nil, customerrors.ErrForbidden
	}
	match, err := utils.VerifyPassword(user.Password, u.Password)
	if err != nil {
		return nil, customerrors.ErrInternalError
	}
	if !match {
		return nil, customerrors.ErrInvalidLogin //401?
	}
	session, sessionErr := s.sessionRepo.CreateSession(cx, user.ID)
	if sessionErr != nil {
		return nil, sessionErr
	}
	return session, nil
}

// deleting the session will force the user to the login page
// and the middlewere will prefent the user to go back without login session
func (s *UserService) LogoutService(cx context.Context, sessionId string) error {
	return s.sessionRepo.DeleteSession(cx, sessionId)
}

// this check if the profided username/email are available
func (s *UserService) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	return s.repo.CheckIfAvailable(cx, col)
}
