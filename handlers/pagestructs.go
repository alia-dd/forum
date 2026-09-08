package handlers

import "gitea.kood.tech/jyrkikarhunen/forum/models"

//todo: add user info
type MainPage struct {
	Posts      []*models.PostView
	Categories []models.Category
}
