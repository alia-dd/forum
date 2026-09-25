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

	http.Redirect(w, r, "/", http.StatusSeeOther)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	//admin validation
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return
	}

	err = h.catRep.DeleteCategory(r.Context(), id)
	if err != nil {
		handleError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
