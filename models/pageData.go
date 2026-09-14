package models

type PageData struct {
	User        *UserInfo
	IsOwner     bool
	PageContent any
	Error       string
}
