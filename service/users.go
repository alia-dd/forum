package service

import (
	"context"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUserService(cx context.Context, u models.UserRegister) error {
	hashedPass, hashErr := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	u.Password = string(hashedPass)
	return s.repo.CreatUser(cx, u)
}

func (s *UserService) GetUserByIdService(cx context.Context, id int) error {
	return s.repo.GetUserById(cx, id)
}

func (s *UserService) CheckIfAvailable(cx context.Context, col string) (bool, error) {
	return s.repo.CheckIfAvailable(cx, col)
}
