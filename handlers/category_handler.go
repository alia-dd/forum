package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
)

type CategoryHandler struct {
	catRep repository.CategoryRepository
}

func NewCategoryHandler(catRep repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{catRep: catRep}
}

func (h *CategoryHandler) NewCategoryForm(w http.ResponseWriter, r *http.Request) {
	//require admin?
	//execute new category template
	fmt.Fprintf(w, "requesting category form")
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	form := CategoryForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}

	if !form.Validate() {
		//execute same template with form values (prefilled) and form errors next to relevant sections
		return
	}

	_, err := h.catRep.CreateCategory(r.Context(), form.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrDuplicateEntry) {
			form.Errors["Name"] = "That category already exists"
			//execute same template with form values (prefilled) and form errors next to relevant sections
			return
		}
		handleError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CategoryHandler) EditCategoryForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	cat, err := h.catRep.GetCategoryByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	//admin validation

	fmt.Fprintf(w, "cat=%v", cat)

	//execute categoryedittemplate with user, cat
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	form := CategoryForm{
		Name: strings.TrimSpace(r.FormValue("name")),
	}

	if !form.Validate() {
		//execute same template with form values (prefilled) and form errors next to relevant sections
		return
	}

	err = h.catRep.UpdateCategory(r.Context(), id, form.Name)
	if err != nil {
		if errors.Is(err, customerrors.ErrDuplicateEntry) {
			form.Errors["Name"] = "That category already exists"
			//execute same template with form values (prefilled) and form errors next to relevant sections
			return
		}
		handleError(w, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	//admin validation
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	err = h.catRep.DeleteCategory(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
