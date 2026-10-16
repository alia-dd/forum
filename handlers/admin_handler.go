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

	usersData, pg, err := paginate(r, func(limit, offset int) ([]models.AdminUserInfo, error) {
		return h.searchRep.SearchAdminUsers(cx, q, limit, offset)
	})
	if err != nil {
		handleError(w, r, err)
		return
	}

	if redirectIfPageTooHigh(w, r, pg, len(usersData)) {
		return
	}

	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "User Table",
			SearchPath: "/admin/users",
			Query:      q,
			Users:      usersData,
			Pagination: pg,
		},
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
	userId, userErr := strconv.Atoi(r.URL.Query().Get("user_id"))
	action := r.URL.Query().Get("action")
	actionType := r.FormValue("ActionType")

	if idErr != nil || userErr != nil || (action != "role" && action != "block" && action != "delete" && actionType != "userAction" && actionType != "postAction" && actionType != "categoryAction") {
		fmt.Println(r.URL.Query().Get("id"), r.URL.Query().Get("user_id"), r.URL.Query().Get("action"), r.URL.Query().Get("ActionType"))
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}
	utils.RenderPartial(w, http.StatusOK, "confirm_secret",
		struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      string
		}{id, userId, action, actionType, ""})
}
func confirmFormError(w http.ResponseWriter, id, userId int, action, actionType, errormgs string) {

	w.Header().Set("HX-Retarget", "#cfm_diologe")
	w.Header().Set("HX-Reswap", "innerHTML")
	utils.RenderPartial(w, http.StatusOK, "confirm_secret",
		struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      string
		}{id, userId, action, actionType, errormgs})
}

func (h *UseHandler) PostAdminUserAction(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}

	cfm_secter := r.FormValue("cfm_secter")
	userId, _ := strconv.Atoi(r.FormValue("user_id"))
	action := r.FormValue("action")
	actionType := r.FormValue("ActionType")
	if userId == user.ID {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	if _, cfmErr := confirmaAdminsecret(cfm_secter); cfmErr != nil {
		confirmFormError(w, 0, userId, action, actionType, "Incorrect confermation secret")
		return
	}

	updateErr := h.service.UpdateUserRoleService(cx, action, userId)

	usersData, err := h.service.GetAllUsersService(cx)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Users: usersData,
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

	posts, pg, err := paginate(r, func(limit, offset int) ([]*models.PostView, error) {
		if q != "" {
			return h.searchRep.SearchAdminPosts(cx, createFTSQuery(q), limit, offset)
		}
		return h.postRep.GetPost(cx, models.PostFilter{Limit: limit, Offset: offset}, -1)
	})
	if err != nil {
		handleError(w, r, err)
		return
	}

	if redirectIfPageTooHigh(w, r, pg, len(posts)) {
		return
	}

	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "Post Table",
			SearchPath: "/admin/posts",
			Query:      q,
			Posts:      posts,
			Pagination: pg,
		},
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
	actionType := r.FormValue("ActionType")
	// check here

	if _, cfmErr := confirmaAdminsecret(cfm_secter); cfmErr != nil {
		confirmFormError(w, postId, userId, action, actionType, "Incorrect confermation secret")
		return
	}
	role := ""
	if ok {
		role = user.Role
	}
	deleteErr := h.postRep.DeletePost(cx, postId, userId, role)

	posts, err := h.postRep.GetPost(cx, models.PostFilter{}, -1)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Posts: posts,
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
