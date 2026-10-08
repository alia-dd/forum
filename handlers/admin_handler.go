package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func (h *UseHandler) GetAdminUserTable(w http.ResponseWriter, r *http.Request) {

	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}

	q := r.URL.Query().Get("q")

	usersData, err := h.searchRep.SearchAdminUsers(cx, q)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "User Table",
			SearchPath: "/admin/users",
			Query:      q,
			Users:      usersData,
		},
	}
	if err != nil {
		fmt.Println(err)
		handleError(w, r, err)
		return
	}
	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}

func (h *UseHandler) GetConfirmUserAction(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}

	id, idErr := strconv.Atoi(r.URL.Query().Get("id"))
	action := r.URL.Query().Get("action")
	endPoint := r.URL.Query().Get("ActionType")

	if idErr != nil || (action != "role" && action != "block" && action != "delete") {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}
	utils.RenderPartial(w, http.StatusOK, "confirm_secret",
		struct {
			ID       int
			Action   string
			EndPoint string
			Error    error
		}{id, action, endPoint, nil})
}

func (h *UseHandler) PostAdminUserAction(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	fmt.Println(">1")
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		fmt.Println(">2")
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	cfm_secter := r.FormValue("cfm_secter")
	userId, _ := strconv.Atoi(r.FormValue("id"))
	action := r.FormValue("action")
	// check here
	fmt.Println(">>", userId, action, cfm_secter)

	if cfm_secter != "12345" {
		utils.RenderPartial(w, http.StatusOK, "confirm_secret", struct {
			ID     int
			Action string
			Error  string
		}{userId, action, "Incorrect confermation secret"})
		return
	}

	updateErr := h.service.UpdateUserRoleService(cx, action, userId)

	usersData, err := h.service.GetAllUsersService(cx)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section: "User Table",
			Users:   usersData,
		},
	}
	if err != nil {
		handleError(w, r, err)
		return

	}
	if updateErr != nil {
		handleError(w, r, updateErr)
		return
	}
	utils.RenderPartial(w, http.StatusOK, "user_table", pageData)
	w.Write([]byte(`<div id="cfm_diologe" hx-swap-oob="true"></div>`))
}

func (h *UseHandler) GetAdminPostTable(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	q := r.URL.Query().Get("q")

	var posts []*models.PostView
	var err error
	if q != "" {
		posts, err = h.searchRep.SearchAdminPosts(cx, createFTSQuery(q))
	} else {
		posts, err = h.postRep.GetPost(cx, models.PostFilter{}, -1)
	}

	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "Post Table",
			SearchPath: "/admin/posts",
			Query:      q,
			Posts:      posts,
		},
	}
	if err != nil {
		fmt.Println(err)
		handleError(w, r, err)
		return
	}

	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}

func (h *UseHandler) PostAdminPostrAction(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	cfm_secter := r.FormValue("cfm_secter")
	postId, _ := strconv.Atoi(r.FormValue("id"))
	// userId, _ := strconv.Atoi(r.FormValue("user_id"))
	action := r.FormValue("action")
	// check here
	fmt.Println(postId, action, cfm_secter)

	if cfm_secter != "12345" {
		utils.RenderPartial(w, http.StatusOK, "confirm_secret", struct {
			ID     int
			Action string
			Error  string
		}{postId, action, "Incorrect confermation secret"})
		return
	}

	// updateErr := h.postRep.DeletePost(cx, postId, userId)

	posts, err := h.postRep.GetPost(cx, models.PostFilter{}, -1)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section: "Post Table",
			Posts:   posts,
		},
	}
	if err != nil {
		handleError(w, r, customerrors.ErrForbidden)
		return

	}
	// if updateErr != nil {
	// 	pageData.Error = updateErr.Error()

	// }
	utils.RenderPartial(w, http.StatusOK, "post_table", pageData)
	w.Write([]byte(`<div id="cfm_diologe" hx-swap-oob="true"></div>`))
}

func (h *UseHandler) GetAdminCommentTable(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	comments, err := h.commentRep.GetCommentsByPost(cx, 0, nil)
	fmt.Println("here")
	fmt.Println(comments, err)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:  "Comment Table",
			Comments: comments,
		},
	}
	if err != nil {
		handleError(w, r, err)
		return
	}

	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}

func (h *UseHandler) GetCategoryTable(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	category, err := h.categoryRep.GetAllCategories(cx)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "Category Table",
			Categories: category,
		},
	}
	if err != nil {
		handleError(w, r, err)
		return
	}

	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}
