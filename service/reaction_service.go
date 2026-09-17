package service

import (
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

type ReactionService struct {
	Reactrepo *repository.ReactionRepository
}

func NewReactionHandler(repo *repository.ReactionRepository) *ReactionService {
	return &ReactionService{Reactrepo: repo}
}
