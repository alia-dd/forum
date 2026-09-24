package models

type PageData[T any] struct {
	User        *UserInfo
	IsOwner     bool
	Query       string
	Scope       string
	PageContent T
	Error       string
}
