package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
)

type ReactionHandler struct {
	service *service.ReactionService
}

func NewReactionHandler(service *service.ReactionService) *ReactionHandler {
	return &ReactionHandler{service: service}
}

func (h *ReactionHandler) Reaction(w http.ResponseWriter, r *http.Request) {
	// cx := r.Context()

	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	targetType := r.FormValue("target_type")
	targetID, _ := strconv.Atoi(r.FormValue("target_id"))
	value, _ := strconv.Atoi(r.FormValue("value"))

	view := models.Reaction{
		Id:         targetID,
		User_id:    user.Id,
		Value:      value,
		TargetType: targetType,
	}
	fmt.Printf("%+v\n", view)
	switch targetType {
	case "post":

	case "comment":

	default:
		handleError(w, r, customerrors.ErrInternalError)
		return
	}
}
