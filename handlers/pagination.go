package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

const pageSize = 25

func pageNumber(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || n < 1 {
		n = 1
	}
	return n
}

func pageOffset(page int) int {
	return (page - 1) * pageSize
}

func createPagination(r *http.Request, page int, hasNext bool) models.Pagination {
	q := r.URL.Query()
	q.Del("page")
	Base := q.Encode()
	if Base != "" {
		Base += "&"
	}

	hasPrev := page > 1

	return models.Pagination{
		Page:     page,
		NextPage: page + 1,
		PrevPage: page - 1,
		HasNext:  hasNext,
		HasPrev:  hasPrev,
		Base:     template.URL(Base),
	}
}

func paginate[T any](r *http.Request, fetch func(limit, offset int) ([]T, error)) ([]T, models.Pagination, error) {
	page := pageNumber(r)

	// try to fetch one extra page's worth of stuff to easily see if there *is* an extra page
	data, err := fetch(pageSize+1, pageOffset(page))
	if err != nil {
		return nil, models.Pagination{}, err
	}

	hasNext := len(data) > pageSize

	// there was an extra page -> slice those results off
	if hasNext {
		data = data[:pageSize]
	}

	return data, createPagination(r, page, hasNext), nil
}

func redirectIfPageTooHigh(w http.ResponseWriter, r *http.Request, pg models.Pagination, count int) bool {
	if pg.Page > 1 && count == 0 {
		http.Redirect(w, r, r.URL.Path+"?"+string(pg.Base)+"page=1", http.StatusSeeOther)
		return true
	}
	return false
}
