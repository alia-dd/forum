package handlers

import (
	"net/http"
	"strings"
	"unicode/utf8"

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
	scope := r.URL.Query().Get("scope")
	if scope != "users" && scope != "posts" && scope != "comments" {
		scope = "all"
	}

	results := SearchResultPage{Query: q, Scope: scope}
	var err error
	if utf8.RuneCountInString(q) >= 3 {
		if scope == "all" || scope == "users" {
			results.Users, err = h.searchRep.SearchUsers(r.Context(), q)
			if err != nil {
				handleError(w, r, err)
				return
			}
		}
		if scope == "all" || scope == "posts" {
			results.Posts, err = h.searchRep.SearchPosts(r.Context(), q)
			if err != nil {
				handleError(w, r, err)
				return
			}
		}
		if scope == "all" || scope == "comments" {
			results.Comments, err = h.searchRep.SearchComments(r.Context(), q)
			if err != nil {
				handleError(w, r, err)
				return
			}
		}
	}

	pageData := models.PageData[SearchResultPage]{
		User: user,
		PageContent: results,
	}

	utils.RenderTemplate(w, http.StatusOK, "search", pageData)
}
