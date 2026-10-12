package server

import (
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

type routeDependencies struct {
	sessionRepo *repository.SessionRepository
	userRepo    *repository.UserRepository

	searchHandler   *handlers.SearchHandler
	reactHandler    *handlers.ReactionHandler
	userHandler     *handlers.UseHandler
	postHandler     *handlers.PostHandler
	categoryHandler *handlers.CategoryHandler
	commentHandler  *handlers.CommentHandler
}
