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
	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}

func (h *UseHandler) GetConfirmUserAction(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	fmt.Println(r.URL.Query().Get("id"), r.URL.Query().Get("user_id"), r.URL.Query().Get("action"), r.URL.Query().Get("ActionType"))

	id, idErr := strconv.Atoi(r.URL.Query().Get("id"))
	userId, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	action := r.URL.Query().Get("action")
	actionType := r.URL.Query().Get("ActionType")
	if idErr != nil || (action != "role" && action != "block" && action != "delete") {

		fmt.Println("idErr", idErr)
		// fmt.Println("userIdErr", userIdErr)
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}
	utils.RenderPartial(w, http.StatusOK, "confirm_secret",
		struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      error
		}{id, userId, action, actionType, nil})
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
	userId, _ := strconv.Atoi(r.FormValue("user_id"))
	action := r.FormValue("action")
	actionType := r.URL.Query().Get("ActionType")
	// check here
	fmt.Println(">>", userId, action, cfm_secter)

	if cfm_secter != "12345" {
		utils.RenderPartial(w, http.StatusOK, "confirm_secret", struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      string
		}{-1, userId, action, actionType, "Incorrect confermation secret"})
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
	posts, err := h.postRep.GetPost(cx, models.PostFilter{}, -1)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section: "Post Table",
			Posts:   posts,
		},
	}
	if err != nil {
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
	userId, _ := strconv.Atoi(r.FormValue("user_id"))
	action := r.FormValue("action")
	actionType := r.URL.Query().Get("ActionType")
	// check here
	fmt.Println(postId, userId, action, cfm_secter)

	if cfm_secter != "12345" {
		utils.RenderPartial(w, http.StatusOK, "confirm_secret", struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      string
		}{-1, userId, action, actionType, "Incorrect confermation secret"})
		return
	}

	deleteErr := h.postRep.DeletePost(cx, postId, userId)

	posts, err := h.postRep.GetPost(cx, models.PostFilter{}, -1)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section: "Post Table",
			Posts:   posts,
		},
	}
	if err != nil {
		handleError(w, r, err)
		return

	}
	if deleteErr != nil {
		handleError(w, r, deleteErr)
		return

	}
	utils.RenderPartial(w, http.StatusOK, "post_table", pageData)
	w.Write([]byte(`<div id="cfm_diologe" hx-swap-oob="true"></div>`))
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
