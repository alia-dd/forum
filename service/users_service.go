package service

import (
	"context"
	"fmt"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"golang.org/x/crypto/bcrypt"
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
	hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return customerrors.ErrInternalError
	}

	u.Password = string(hashedPass)
	return s.repo.RegisterUser(cx, u)
}

// register new use
func (s *UserService) GetUserService(cx context.Context, username string) (models.UserInfo, error) {

	return s.repo.FetchUserDataByUserName(cx, username)
}

// patch user
func (s *UserService) UpdateUserService(cx context.Context, u models.UserUpdate) error {

	if validationsErr := u.Isvalid(); validationsErr != nil {
		return validationsErr
	}
	if u.Username != nil {
		notAvailable, availableErr := s.repo.CheckIfAvailableExcludeUser(cx, *u.Username, u.Id)
		if availableErr != nil {
			return availableErr
		}
		if notAvailable {
			return customerrors.ErrDuplicateEntry
		}
	}

	if u.Email != nil {
		notAvailable, availableErr := s.repo.CheckIfAvailableExcludeUser(cx, *u.Email, u.Id)
		if availableErr != nil {
			return availableErr
		}
		if notAvailable {
			return customerrors.ErrDuplicateEntry
		}
	}
	return s.repo.UpdateUserData(cx, u)
}

func (s *UserService) ChangePasswordService(cx context.Context, userID int, currentPass, newPass string) error {
	currentPassHash, fetchErr := s.repo.GetPasswordHash(cx, userID)
	if fetchErr != nil {
		return fetchErr
	}
	if compareErr := bcrypt.CompareHashAndPassword([]byte(currentPassHash), []byte(currentPass)); compareErr != nil {
		fmt.Println("not a match")
		return customerrors.ErrIncorrectPassword
	}
	newPassHash, hashErr := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
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
	// fmt.Println(user)
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if compareErr != nil {
		return nil, customerrors.ErrInvalidData
	}
	session, sessionErr := s.sessionRepo.CreateSession(cx, user.Id)
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
