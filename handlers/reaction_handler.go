package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type ReactionHandler struct {
	service *service.ReactionService
	postRep repository.PostRepository
}

func NewReactionHandler(service *service.ReactionService, postRep repository.PostRepository) *ReactionHandler {
	return &ReactionHandler{service: service, postRep: postRep}
}

func (h *ReactionHandler) Reaction(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()

	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}

	targetType := r.FormValue("target_type")
	targetID, idErr := strconv.Atoi(r.FormValue("target_id"))
	value, valErr := strconv.Atoi(r.FormValue("value"))
	if idErr != nil || valErr != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	view := models.Reaction{
		Id:         targetID,
		User_id:    user.Id,
		Value:      value,
		TargetType: targetType,
	}
	fmt.Printf("%+v\n", view)

	var pv *models.PostView
	var err error
	switch targetType {
	case "post":
		if err := h.service.SetPostReact(cx, view); err != nil {
			handleError(w, r, err)
			return
		}
		pv, err = h.postRep.GetPostByID(cx, targetID, user.Id)
	case "comment":
		if err := h.service.SetCommentReact(cx, view); err != nil {
			handleError(w, r, err)
			return
		}
	default:
		handleError(w, r, customerrors.ErrInternalError)
		return
	}
	if err != nil {

	}
	utils.RenderPartial(w, http.StatusAccepted, "react", pv)
}
