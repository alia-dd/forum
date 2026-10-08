package handlers

import "gitea.kood.tech/jyrkikarhunen/forum/models"

//todo: add user info
type MainPage struct {
	Section    string
	SearchPath string
	Query      string
	Posts      []*models.PostView
	Categories []models.Category
	Users      []models.AdminUserInfo
	Comments   []*models.CommentView
	Category   []models.Category
}

type PostPage struct {
	Post     *models.PostView
	Comments []*models.CommentView
}

type PostFormPage struct {
	PostID     int
	Form       PostForm
	Categories []models.Category
}

type CategoryFormPage struct {
	CategoryID int
	Form       CategoryForm
}

type SearchResultPage struct {
	Query    string
	Scope    string
	Sort     string
	Users    []models.UserResult
	Posts    []models.PostResult
	Comments []models.CommentResult
}
