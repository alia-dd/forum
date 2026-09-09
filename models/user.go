package models

import "time"

type UserRegister struct {
	Id              int       `json:"user_id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
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
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}
