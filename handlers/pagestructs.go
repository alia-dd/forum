package handlers

import "gitea.kood.tech/jyrkikarhunen/forum/models"

//todo: add user info
type MainPage struct {
	Posts      []*models.PostView
	Categories []models.Category
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
	Users    []models.UserResult
	Posts    []models.PostResult
	Comments []models.CommentResult
}
