package service

import (
	"context"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

type ReactionService struct {
	repo *repository.ReactionRepository
}

func NewReactionService(repo *repository.ReactionRepository) *ReactionService {
	return &ReactionService{repo: repo}
}

func (rec *ReactionService) SetPostReact(cx context.Context, react models.Reaction) {

}

func (rec *ReactionService) SetCommentReact(cx context.Context, react models.Reaction) {

}
