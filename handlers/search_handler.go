package handlers

import (
	"net/http"
	"strings"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

const searchAllCap = 10 // applied for each separate search category if searching for "all"

type SearchHandler struct {
	searchRep *repository.SearchRepository
}

func NewSearchHandler(s *repository.SearchRepository) *SearchHandler {
	return &SearchHandler{searchRep: s}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user_session").(*models.UserInfo)
	ctx := r.Context()

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

	if scope == "all" {
		results.Users, err = h.searchRep.SearchUsers(ctx, q, searchAllCap, 0)
		if err != nil {
			handleError(w, r, err)
			return
		}

		results.Posts, err = h.searchRep.SearchPosts(ctx, ftsQuery, sort, searchAllCap, 0)
		if err != nil {
			handleError(w, r, err)
			return
		}

		results.Comments, err = h.searchRep.SearchComments(ctx, ftsQuery, sort, searchAllCap, 0)
		if err != nil {
			handleError(w, r, err)
			return
		}
	}

	// user search only uses the basic sql query because it needs to support partial wildcard matches
	if scope == "users" {
		results.Users, results.Pagination, err = paginate(r, func(limit, offset int) ([]models.UserResult, error) {
			return h.searchRep.SearchUsers(ctx, q, limit, offset)
		})
		if err != nil {
			handleError(w, r, err)
			return
		}

		if redirectIfPageTooHigh(w, r, results.Pagination, len(results.Users)) {
			return
		}
	}
	if scope == "posts" {
		results.Posts, results.Pagination, err = paginate(r, func(limit, offset int) ([]models.PostResult, error) {
			return h.searchRep.SearchPosts(ctx, ftsQuery, sort, limit, offset)
		})
		if err != nil {
			handleError(w, r, err)
			return
		}

		if redirectIfPageTooHigh(w, r, results.Pagination, len(results.Posts)) {
			return
		}
	}
	if scope == "comments" {
		results.Comments, results.Pagination, err = paginate(r, func(limit, offset int) ([]models.CommentResult, error) {
			return h.searchRep.SearchComments(ctx, ftsQuery, sort, limit, offset)
		})
		if err != nil {
			handleError(w, r, err)
			return
		}

		if redirectIfPageTooHigh(w, r, results.Pagination, len(results.Comments)) {
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
