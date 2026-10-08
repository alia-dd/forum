package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type CategoryHandler struct {
	catRep repository.CategoryRepository
}

func NewCategoryHandler(catRep repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{catRep: catRep}
}

func (h *CategoryHandler) NewCategoryForm(w http.ResponseWriter, r *http.Request) {
	//require admin?
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}

	pageData := models.PageData[CategoryFormPage]{
		User:        user,
		PageContent: CategoryFormPage{},
	}
	utils.RenderTemplate(w, http.StatusOK, "category_form", pageData)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	form := CategoryForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}

	if !form.Validate() {
		utils.RenderTemplate(w, http.StatusOK, "category_form", models.PageData[CategoryFormPage]{
			User: user, PageContent: CategoryFormPage{Form: form},
		})
		return
	}

	_, err := h.catRep.CreateCategory(r.Context(), form.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrDuplicateEntry) {
			form.Errors["Name"] = "That category already exists"
			utils.RenderTemplate(w, http.StatusOK, "category_form", models.PageData[CategoryFormPage]{
				User: user, PageContent: CategoryFormPage{Form: form},
			})
			return
		}
		handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) EditCategoryForm(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	cat, err := h.catRep.GetCategoryByID(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}

	pageData := models.PageData[CategoryFormPage]{
		User: user,
		PageContent: CategoryFormPage{
			CategoryID: id,
			Form: CategoryForm{
				Name: cat.Name,
			},
		},
	}

	//admin validation

	utils.RenderTemplate(w, http.StatusOK, "category_form", pageData)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		utils.RedirectTologin(w, r)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	form := CategoryForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}

	if !form.Validate() {
		utils.RenderTemplate(w, http.StatusOK, "category_form", models.PageData[CategoryFormPage]{
			User: user, PageContent: CategoryFormPage{CategoryID: id, Form: form},
		})
		return
	}

	err = h.catRep.UpdateCategory(r.Context(), id, form.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrDuplicateEntry) {
			form.Errors["Name"] = "That category already exists"
			utils.RenderTemplate(w, http.StatusOK, "category_form", models.PageData[CategoryFormPage]{
				User: user, PageContent: CategoryFormPage{CategoryID: id, Form: form},
			})
			return
		}
		handleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		handleError(w, r, customerrors.ErrForbidden)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	cfm_secter := r.FormValue("cfm_secter")
	action := r.FormValue("action")
	actionType := r.URL.Query().Get("ActionType")
	fmt.Println(">", id)
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}
	if cfm_secter != "12345" {
		utils.RenderPartial(w, http.StatusOK, "confirm_secret", struct {
			ID         int
			User_id    int
			Action     string
			ActionType string
			Error      string
		}{id, 0, action, actionType, "Incorrect confermation secret"})
		return
	}
	err = h.catRep.DeleteCategory(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	categories, fetchErr := h.catRep.GetAllCategories(r.Context())
	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Section:    "Post Table",
			Categories: categories,
		},
	}
	if fetchErr != nil {
		handleError(w, r, fetchErr)
		return

	}
	utils.RenderPartial(w, http.StatusOK, "category_table", pageData)
	w.Write([]byte(`<div id="cfm_diologe" hx-swap-oob="true"></div>`))
}

// func (h *UseHandler) GetAdminUserTable(w http.ResponseWriter, r *http.Request) {

// func (h *UseHandler) GetConfirmUserAction(w http.ResponseWriter, r *http.Request) {
// 	user, ok := r.Context().Value("user_session").(*models.UserInfo)
// 	if !ok || user.Role != "admin" {
// 		handleError(w, r, customerrors.ErrForbidden)
// 		return
// 	}
// 	fmt.Println(r.URL.Query().Get("id"), r.URL.Query().Get("user_id"), r.URL.Query().Get("action"), r.URL.Query().Get("ActionType"))

// 	id, idErr := strconv.Atoi(r.URL.Query().Get("id"))
// 	userId, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
// 	action := r.URL.Query().Get("action")
// 	actionType := r.URL.Query().Get("ActionType")
// 	if idErr != nil || (action != "role" && action != "block" && action != "delete") {

// 		fmt.Println("idErr", idErr)
// 		handleError(w, r, customerrors.ErrBadRequest)
// 		return
// 	}
// 	utils.RenderPartial(w, http.StatusOK, "confirm_secret",
// 		struct {
// 			ID         int
// 			User_id    int
// 			Action     string
// 			ActionType string
// 			Error      error
// 		}{id, userId, action, actionType, nil})
// }
