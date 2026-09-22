package models

type PageData[T any] struct {
	User        *UserInfo
	IsOwner     bool
	PageContent T
	Error       string
}
