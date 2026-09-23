package handlers

import (
	"net/http"
	"strconv"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type CommentHandler struct {
	comRep repository.CommentRepository
}

func NewCommentHandler(comRep repository.CommentRepository) *CommentHandler {
	return &CommentHandler{comRep: comRep}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("userComment"))
	if content == "" {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	var parentCommentID *int
	parentIDStr := r.FormValue("parent_id")
	if parentIDStr != "" {
		if id, err := strconv.Atoi(parentIDStr); err == nil {
			parentCommentID = &id
		}
	}

	input := models.CommentInput{
		UserID:          user.Id,
		ParentPostID:    postID,
		ParentCommentID: parentCommentID,
		Content:         content,
	}

	comID, err := h.comRep.CreateComment(r.Context(), input)
	if err != nil {
		handleError(w, r, err)
		return
	}

	newComment, err := h.comRep.GetCommentByID(r.Context(), comID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	newComment.IsOwner = (user.Id == newComment.UserID)

	newCount, err := h.comRep.CountByPost(r.Context(), postID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	responsePayload := struct {
		NewComment *models.CommentView
		NewCount   int
	}{
		NewComment: newComment,
		NewCount:   newCount,
	}

	utils.RenderPartial(w, http.StatusOK, "comment_response", responsePayload)
}

func (h *CommentHandler) GetCommentEditForm(w http.ResponseWriter, r *http.Request) {
	commentID, _ := parseID(r, "commentID")

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	utils.RenderPartial(w, http.StatusOK, "comment_card_edit", comment)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	commentID, _ := parseID(r, "commentID")
	newContent := r.FormValue("content")

	err := h.comRep.UpdateCommentContent(r.Context(), commentID, user.Id, newContent)
	if err != nil {
		handleError(w, r, err)
		return
	}

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	comment.IsOwner = (user.Id == comment.UserID)

	utils.RenderPartial(w, http.StatusOK, "comment_card", comment)
}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	commentID, _ := parseID(r, "commentID")

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID)
	if err != nil {
		handleError(w, r, err)
		return
	}
	comment.IsOwner = (user.Id == comment.UserID)

	utils.RenderPartial(w, http.StatusOK, "comment_card", comment)
}
