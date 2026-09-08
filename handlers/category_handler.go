package handlers

import (
	"errors"
	"fmt"
	"net/http"
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
