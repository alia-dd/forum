package handlers

import (
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

type ReactionHandler struct {
	repo *repository.ReactionRepository
}

func NewReactionHandler(repo *repository.ReactionRepository) *ReactionHandler {
	return &ReactionHandler{repo: repo}
}

func (h *UseHandler) PostReaction(w http.ResponseWriter, r *http.Request) {
	// cx := r.Context()

	// _, ok := r.Context().Value("user_session").(*models.UserInfo)
	// if !ok {
	// 	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	// 	return
	// }
	// value := r.PathValue("reactionValue")

	// if PostErr := h.service.CreateUserService(cx, userData); PostErr != nil {
	// 	handleError(w, customerrors.ErrBadRequest)
	// 	return
	// }

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
