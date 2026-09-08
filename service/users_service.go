package service

import (
	"context"

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

	hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return customerrors.ErrInternalError
	}

	u.Password = string(hashedPass)
	return s.repo.RegisterUser(cx, u)
}

// authernticate use with their email/usename and password
func (s *UserService) AuthenticateUserService(cx context.Context, u models.UserLogin) (*models.Session, error) {

	user, fetchErr := s.repo.AuthenticateUser(cx, u.Username)
	if fetchErr != nil {
		return nil, fetchErr
	}
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

func (s *UserService) DeleteUserSession(cx context.Context, sessionId string) error {
	return s.sessionRepo.DeleteSession(cx, sessionId)
}

func (s *UserService) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	return s.repo.CheckIfAvailable(cx, col)
}
