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
		HandleError(w, r, customerrors.ErrBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		HandleError(w, r, customerrors.ErrBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("userComment"))
	if content == "" {
		HandleError(w, r, customerrors.ErrBadRequest)
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
		UserID:          user.ID,
		ParentPostID:    postID,
		ParentCommentID: parentCommentID,
		Content:         content,
	}

	comID, err := h.comRep.CreateComment(r.Context(), input)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	newComment, err := h.comRep.GetCommentByID(r.Context(), comID, &user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	newComment.IsOwner = (user.ID == newComment.UserID)

	newCount, err := h.comRep.CountByPost(r.Context(), postID)
	if err != nil {
		HandleError(w, r, err)
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
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	commentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID, &user.ID)
	if err != nil {
		HandleError(w, r, err)
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

	err := h.comRep.UpdateCommentContent(r.Context(), commentID, user.ID, user.Role, newContent)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID, &user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	comment.IsOwner = (user.ID == comment.UserID)
	comment.UserRole = user.Role
	utils.RenderPartial(w, http.StatusOK, "comment_card", comment)
}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	commentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	comment, err := h.comRep.GetCommentByID(r.Context(), commentID, &user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}
	comment.IsOwner = (user.ID == comment.UserID)
	comment.UserRole = user.Role
	utils.RenderPartial(w, http.StatusOK, "comment_card", comment)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	commentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	postID, err := parseID(r, "postID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	// cascade is set to false only admin can delete comment and all of it's repliec
	err = h.comRep.DeleteComment(r.Context(), commentID, user.ID, user.Role, false)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	newCount, err := h.comRep.CountByPost(r.Context(), postID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	responsePayload := struct {
		NewCount int
	}{
		NewCount: newCount,
	}

	utils.RenderPartial(w, http.StatusOK, "comment_delete_response", responsePayload)
}

func (h *CommentHandler) GetCommentReplies(w http.ResponseWriter, r *http.Request) {
	var userID *int
	role := ""
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if ok {
		userID = &user.ID
		role = user.Role
	}

	commentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	replies, err := h.comRep.GetRepliesByCommentID(r.Context(), commentID, userID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	for _, reply := range replies {
		reply.UserRole = role
	}

	utils.RenderPartial(w, http.StatusOK, "comment_replies", replies)
}

func (h *CommentHandler) GetReplyForm(w http.ResponseWriter, r *http.Request) {
	postID, err := parseID(r, "postID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	commentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	data := map[string]any{
		"ParentPostID": postID,
		"ParentID":     commentID,
	}

	utils.RenderPartial(w, http.StatusOK, "comment_reply", data)
}

func (h *CommentHandler) CreateReply(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	postID, err := parseID(r, "postID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	parentCommentID, err := parseID(r, "commentID")
	if err != nil {
		HandleError(w, r, err)
		return
	}

	content := r.FormValue("content")

	input := models.CommentInput{
		UserID:          user.ID,
		ParentPostID:    postID,
		ParentCommentID: &parentCommentID,
		Content:         content,
	}

	replyID, err := h.comRep.CreateComment(r.Context(), input)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	replyView, err := h.comRep.GetCommentByID(r.Context(), replyID, &user.ID)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	utils.RenderPartial(w, http.StatusOK, "comment_card", replyView)
}
