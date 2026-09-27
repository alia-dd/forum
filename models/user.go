package models

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
)

type UserRegister struct {
	ID              int       `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Image           string    `json:"imagePath"`
	Bio             string    `json:"bio"`
	Password        string    `json:"password"`
	ConformPassword string    `json:"ConformPassword"`
	CreatedAt       time.Time `json:"createdat"`
	UpdatedAt       time.Time `json:"updatedat"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Name      string    `json:"name,omitempty"`
	Image     string    `json:"imagePath,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type PublicUserInfo struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Name      string    `json:"name,omitempty"`
	Image     string    `json:"imagePath,omitempty"`
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"createdat"`
}

type UserUpdate struct {
	ID       int     `json:"id"`
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Image    *string `json:"imagePath"`
	Bio      *string `json:"bio"`
}

func (u *UserRegister) Isvalid() error {
	if u.Username != "" && !IsValidName(u.Username) {
		return customerrors.ErrInvalidName
	}
	if u.Email != "" && !IsValidEmail(u.Email) {
		return customerrors.ErrInvalidData
	}
	if u.Password != "" && u.ConformPassword != "" && (u.Password != u.ConformPassword) {
		return customerrors.ErrIncorrectPassword
	}
	return nil
}

func (u *UserUpdate) Isvalid() error {
	if u.Username != nil && !IsValidName(*u.Username) {
		return customerrors.ErrInvalidName
	}
	if u.Email != nil && !IsValidEmail(*u.Email) {
		fmt.Println(u.Email)
		return customerrors.ErrInvalidData
	}
	if u.Bio != nil && len(*u.Bio) > 200 {
		return customerrors.ErrInvalidData
	}
	return nil
}
func IsValidName(name string) bool {
	r, _ := regexp.Compile("^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$")
	if r.MatchString(name) {
		return true
	}
	fmt.Println("not a march")
	return false
}

func IsValidEmail(email string) bool {
	validEmail, emailErr := mail.ParseAddress(email)
	return emailErr == nil && validEmail.Address == email
}
