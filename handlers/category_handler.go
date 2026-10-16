package handlers

import (
	"errors"
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
		HandleError(w, r, customerrors.ErrBadRequest)
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
		HandleError(w, r, err)
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
		HandleError(w, r, customerrors.ErrBadRequest)
		return
	}

	cat, err := h.catRep.GetCategoryByID(r.Context(), id)
	if err != nil {
		HandleError(w, r, err)
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
		HandleError(w, r, customerrors.ErrBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		HandleError(w, r, customerrors.ErrBadRequest)
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
		HandleError(w, r, err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok || user.Role != "admin" {
		HandleError(w, r, customerrors.ErrForbidden)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	cfm_sector := r.FormValue("cfm_sector")
	action := r.FormValue("action")
	actionType := r.FormValue("ActionType")

	if err != nil {
		HandleError(w, r, customerrors.ErrBadRequest)
		return
	}
	if cfm_sector != "12345" {
		confirmFormError(w, id, 0, action, actionType, "Incorrect confermation secret")
		return
	}
	err = h.catRep.DeleteCategory(r.Context(), id)
	if err != nil {
		HandleError(w, r, err)
		return
	}

	categories, fetchErr := h.catRep.GetAllCategories(r.Context())

	pageData := models.PageData[MainPage]{
		User: user,
		PageContent: MainPage{
			Categories: categories,
		},
	}
	if fetchErr != nil {
		HandleError(w, r, fetchErr)
		return

	}

	utils.RenderPartial(w, http.StatusOK, "category_table", pageData)
	w.Write([]byte(`<div id="cfm_dialogue" hx-swap-oob="true"></div>`))
}
