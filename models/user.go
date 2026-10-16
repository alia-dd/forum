package models

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
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
	ConfirmPassword string    `json:"ConfirmPassword"`
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
	Role      string    `json:"role,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type AdminUserInfo struct {
	Id        int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      int       `json:"role"`
	Image     string    `json:"imagePath"`
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
	DeletedAt bool
}
type ProfilePage struct {
	Profile    PublicUserInfo
	Tab        string
	Posts      []*PostView
	Comments   []*CommentView
	LikedPosts []*PostView
	Pagination Pagination
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
	if !IsValidName(u.Username) {
		return customerrors.ErrInvalidName
	}
	if !IsValidEmail(u.Email) {
		return customerrors.ErrInvalidEmail
	}
	if u.Password == "" || u.ConfirmPassword == "" || (u.Password != u.ConfirmPassword) {
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
		return customerrors.ErrInvalidEmail
	}
	if u.Bio != nil && len(*u.Bio) > 500 {
		return customerrors.ErrInvalidData
	}
	return nil
}

func IsValidName(name string) bool {
	r, _ := regexp.Compile("^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$")
	if r.MatchString(name) && len(name) >= 3 && len(name) <= 20 {
		return true
	}
	fmt.Println("not a march")
	return false
}

func IsValidEmail(email string) bool {
	validEmail, emailErr := mail.ParseAddress(email)
	correct := emailErr == nil && validEmail.Address == email
	res := strings.Split(email, "@")
	if len(res) != 2 {
		return false
	}
	f := len(res[0])
	s := len(res[1])
	extraCheck := f > 0 && f <= 64 && s > 0 && s <= 255-f
	return correct && extraCheck
}
