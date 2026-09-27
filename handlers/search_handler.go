package handlers

import (
	"net/http"
	"strings"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type SearchHandler struct {
	searchRep *repository.SearchRepository
}

func NewSearchHandler(s *repository.SearchRepository) *SearchHandler {
	return &SearchHandler{searchRep: s}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user_session").(*models.UserInfo)

	q := strings.TrimSpace(r.URL.Query().Get("q"))

	if q == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	scope := r.URL.Query().Get("scope")
	if scope != "users" && scope != "posts" && scope != "comments" {
		scope = "all"
	}

	sort := r.URL.Query().Get("sort")
	if sort != "relevance" && sort != "time" {
		sort = "relevance"
	}

	results := SearchResultPage{Query: q, Scope: scope, Sort: sort}
	var err error
	if !validSearchLength(q) {
		pageData := models.PageData[SearchResultPage]{
			User:        user,
			Query:       q,
			Scope:       scope,
			PageContent: results,
			Error:       "Search must be at least 3 characters.",
		}

		utils.RenderTemplate(w, http.StatusBadRequest, "search", pageData)
		return
	}

	// wrap every word in quotes because dealing with the full potential power of FTS5 queries is not worth it here
	ftsQuery := createFTSQuery(q)

	// user search only uses the basic sql query because it needs to support partial wildcard matches
	if scope == "all" || scope == "users" {
		results.Users, err = h.searchRep.SearchUsers(r.Context(), q)
		if err != nil {
			handleError(w, r, err)
			return
		}
	}
	if scope == "all" || scope == "posts" {
		results.Posts, err = h.searchRep.SearchPosts(r.Context(), ftsQuery, sort)
		if err != nil {
			handleError(w, r, err)
			return
		}
	}
	if scope == "all" || scope == "comments" {
		results.Comments, err = h.searchRep.SearchComments(r.Context(), ftsQuery, sort)
		if err != nil {
			handleError(w, r, err)
			return
		}
	}

	pageData := models.PageData[SearchResultPage]{
		User:        user,
		Query:       q,
		Scope:       scope,
		PageContent: results,
	}

	utils.RenderTemplate(w, http.StatusOK, "search", pageData)
}
