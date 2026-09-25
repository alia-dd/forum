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

func (rec *ReactionService) SetPostReact(cx context.Context, react models.Reaction) error {

	val, fetchErr := rec.repo.GetPostReaction(cx, react)
	if fetchErr != nil {
		return fetchErr
	}
	if val == nil {
		return rec.repo.CreatePostReaction(cx, react)
	}
	if *val == react.Value {
		return rec.repo.DeletePostReaction(cx, react)
	} else {
		return rec.repo.UpdatePostReaction(cx, react)
	}
}

func (rec *ReactionService) SetCommentReact(cx context.Context, react models.Reaction) error {
	val, fetchErr := rec.repo.GetCommentReaction(cx, react)
	if fetchErr != nil {
		return fetchErr
	}
	if val == nil {
		return rec.repo.CreateCommentReaction(cx, react)
	}
	if *val == react.Value {
		return rec.repo.DeleteCommentReaction(cx, react)
	} else {
		return rec.repo.UpdateCommentReaction(cx, react)
	}
}
