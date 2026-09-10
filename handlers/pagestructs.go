package handlers

import "gitea.kood.tech/jyrkikarhunen/forum/models"

//todo: add user info
type MainPage struct {
	Posts      []*models.PostView
	Categories []models.Category
}

type PostPage struct {
	Post *models.PostView
	//comments
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
