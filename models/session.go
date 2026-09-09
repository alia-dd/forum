package models

import "time"

type Session struct {
	SesssionId string    `json:"sesssionId"`
	UserID     int       `json:"userID"`
	ExpiresAt  time.Time `json:"expiresAt"`
	Createdat  time.Time `json:"createdat"`
}
