package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func (h *UseHandler) GetAdminUserTable(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}
	usersData, err := h.service.GetAllUsersService(cx)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Users: usersData,
		},
	}
	if err != nil {
		pageData.Error = err.Error()
	}
	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}

func (h *UseHandler) GetAdminUserAction(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}
	userId, _ := strconv.Atoi(r.FormValue("id"))
	action := r.FormValue("action")
	fmt.Println(userId, action)
	updateErr := h.service.UpdateUserRoleService(cx, action, userId)
	fmt.Println(">", updateErr)
	usersData, err := h.service.GetAllUsersService(cx)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Users: usersData,
		},
	}
	if err != nil {
		pageData.Error = err.Error()

	}
	if updateErr != nil {
		pageData.Error = updateErr.Error()

	}

	// utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
	utils.RenderPartial(w, http.StatusOK, "user_table", pageData)
}
func (h *UseHandler) GetAdminPostTable(w http.ResponseWriter, r *http.Request) {
	cx := r.Context()
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}
	posts, err := h.postRep.GetPost(cx, models.PostFilter{}, 0)
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Posts: posts,
		},
	}
	if err != nil {
		pageData.Error = err.Error()
	}

	utils.RenderTemplate(w, http.StatusOK, "admin_page", pageData)
}
