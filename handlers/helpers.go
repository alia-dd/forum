package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

// postform
type PostForm struct {
	//userdata
	Title          string
	Content        string
	MainCategoryID int
	CategoryIDs    []int
	Errors         map[string]string // [field]message
}

// maxlen of title and content
func (f *PostForm) Validate() bool {
	f.Errors = map[string]string{}

	title := strings.TrimSpace(f.Title)
	content := strings.TrimSpace(f.Content)

	if title == "" {
		f.Errors["Title"] = "Title is required"
	} else if utf8.RuneCountInString(title) > 200 {
		f.Errors["Title"] = "Title is too long"
	}
	if content == "" {
		f.Errors["Content"] = "Content is required"
	} else if utf8.RuneCountInString(content) > 10000 {
		f.Errors["Content"] = "Content is too long"
	}
	if f.MainCategoryID == 0 {
		f.Errors["MainCategoryID"] = "Main category is required"
	}
	return len(f.Errors) == 0
}

func (f PostForm) HasCategory(id int) bool {
	for _, cat := range f.CategoryIDs {
		if cat == id {
			return true
		}
	}
	return false
}

// post helpers
func (h *PostHandler) parsePostFilter(ctx context.Context, r *http.Request) (models.PostFilter, error) {
	var pf models.PostFilter
	q := r.URL.Query()

	cat := q.Get("category")
	if cat != "" {
		id, err := strconv.Atoi(cat)
		if err != nil {
			return pf, err
		}
		pf.CategoryID = &id
	}

	author := q.Get("author")
	if author != "" {
		id, err := h.getAuthorID(ctx, author)
		if err != nil {
			return pf, err
		}
		pf.AuthorID = &id
	}

	liked := q.Get("liked")
	if liked == "true" {
		if user, ok := ctx.Value("user_session").(*models.UserInfo); ok {
			id := user.Id
			pf.LikedByID = &id
		}
	}

	return pf, nil

}

func (h *PostHandler) getAuthorID(ctx context.Context, author string) (int, error) {
	if author == "me" {
		user, ok := ctx.Value("user_session").(*models.UserInfo)
		if !ok {
			return 0, customerrors.ErrBadRequest
		}
		return user.Id, nil
	}

	authorID, err := strconv.Atoi(author)
	if err == nil {
		if authorID <= 0 {
			return 0, customerrors.ErrBadRequest
		}
		return authorID, nil
	}
	//todo: get userid by username and return that or error if not found

	//placeholder return
	return 0, nil
}

func (h *PostHandler) parsePostForm(w http.ResponseWriter, r *http.Request) (PostForm, bool) {
	if err := r.ParseForm(); err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return PostForm{}, false
	}

	mainCat, err := strconv.Atoi(r.FormValue("main_category"))
	if err != nil {
		handleError(w, r, customerrors.ErrBadRequest)
		return PostForm{}, false
	}

	var categoryIDs []int
	for _, cat := range r.Form["category"] {
		id, err := strconv.Atoi(cat)
		if err != nil {
			handleError(w, r, customerrors.ErrBadRequest)
			return PostForm{}, false
		}
		if id != mainCat { // avoid duplicating maincategory into categorylist
			categoryIDs = append(categoryIDs, id)
		}
	}

	return PostForm{
		Title:          strings.TrimSpace(r.FormValue("title")),
		Content:        strings.TrimSpace(r.FormValue("content")),
		MainCategoryID: mainCat,
		CategoryIDs:    categoryIDs,
	}, true
}

// categoryform
type CategoryForm struct {
	//userdata
	Name   string
	Errors map[string]string
}

func (f *CategoryForm) Validate() bool {
	f.Errors = map[string]string{}

	name := strings.TrimSpace(f.Name)

	if name == "" {
		f.Errors["Name"] = "Category name is required"
	} else if utf8.RuneCountInString(name) > 40 {
		f.Errors["Name"] = "Category name is too long"
	}
	return len(f.Errors) == 0
}

func SaveUploadedFile(file multipart.File) (string, error) {
	allowedTypes := []string{"image/jpeg", "image/png"}
	ext, validationErr := validateMIMEType(file, allowedTypes)
	if validationErr != nil {
		return "", validationErr
	}

	uploadDirectory := "static/assets/images"

	if err := os.MkdirAll(uploadDirectory, 0755); err != nil {
		return "", customerrors.ErrInternalError
	}
	saveName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), ext)
	path := filepath.Join(uploadDirectory, saveName)

	dst, err := os.Create(path)
	if err != nil {
		return "", customerrors.ErrInternalError
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {

		return "", customerrors.ErrInternalError
	}

	return saveName, nil
}

func validateMIMEType(file multipart.File, allowedTypes []string) (string, error) {
	// Read the first 512 bytes for detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", customerrors.ErrInternalError
	}

	// Reset file pointer
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", customerrors.ErrInternalError
	}

	contentType := http.DetectContentType(buffer[:n])
	ext := ".jpg"
	if slices.Contains(allowedTypes, contentType) {
		if contentType == "image/png" {
			ext = ".png"
		}
		return ext, nil
	}

	return "", customerrors.ErrBadRequest
}
