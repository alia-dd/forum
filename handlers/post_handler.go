package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

type PostHandler struct {
	postRep repository.PostRepository
	catRep  repository.CategoryRepository
}

func NewPostHandler(postRep repository.PostRepository, catRep repository.CategoryRepository) *PostHandler {
	return &PostHandler{postRep: postRep, catRep: catRep}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user_session").(*models.UserInfo) // guests allowed

	filter, err := h.parsePostFilter(r.Context(), r)
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	posts, err := h.postRep.GetPost(r.Context(), filter)
	if err != nil {
		handleError(w, err)
		return
	}
	ids := make([]int, len(posts))
	for i, post := range posts {
		ids[i] = post.Post.ID
	}
	cats, err := h.catRep.GetCatByPostIDs(r.Context(), ids)
	if err != nil {
		handleError(w, err)
		return
	}
	for _, post := range posts {
		post.Categories = cats[post.Post.ID]
	}
	allCats, err := h.catRep.GetAllCategories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Println("cats: ", allCats)
	fmt.Println("post: ", posts)

	var pageData models.PageData
	pageData = models.PageData{
		User: user,
		PageContent: MainPage{
			Posts:      posts,
			Categories: allCats,
		},
	}

	fmt.Println(pageData)
	utils.RenderTemplate(w, http.StatusOK, "home", pageData)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user_session").(*models.UserInfo) // guests allowed
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	post, err := h.postRep.GetPostByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	cats, err := h.catRep.GetCatByPostIDs(r.Context(), []int{id})
	if err != nil {
		handleError(w, err)
		return
	}
	post.Categories = cats[id]

	//get comments linked to post or its subcomments
	pageData := models.PageData{
		User: user,
		PageContent: PostPage{
			Post: post,
		},
	}

	utils.RenderTemplate(w, http.StatusOK, "post", pageData)
}

func (h *PostHandler) NewPostForm(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	cats, err := h.catRep.GetAllCategories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	pageData := models.PageData{
		User: user,
		PageContent: PostFormPage{
			Categories: cats,
		},
	}

	utils.RenderTemplate(w, http.StatusOK, "post_form", pageData)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	form, ok := h.parsePostForm(w, r)
	if !ok {
		return
	}

	if !form.Validate() {
		cats, err := h.catRep.GetAllCategories(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}
		pageData := models.PageData{
			User: user,
			PageContent: PostFormPage{
				Form:       form,
				Categories: cats,
			},
		}
		utils.RenderTemplate(w, http.StatusOK, "post_form", pageData)
		return
	}

	postInput := models.PostInput{
		UserID:         user.Id,
		Title:          form.Title,
		Content:        form.Content,
		MainCategoryID: form.MainCategoryID,
		CategoryIDs:    form.CategoryIDs,
	}
	id, err := h.postRep.CreatePost(r.Context(), postInput)
	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *PostHandler) EditPostForm(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	post, err := h.postRep.GetPostByID(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	if user.Id != post.Post.UserID {
		handleError(w, customerrors.ErrForbidden)
		return
	}

	cats, err := h.catRep.GetAllCategories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	postCats, err := h.catRep.GetCatByPostIDs(r.Context(), []int{id})
	if err != nil {
		handleError(w, err)
		return
	}

	form := PostForm{Title: post.Post.Title, Content: post.Post.Content}
	for _, c := range postCats[id] {
		if c.IsMain {
			form.MainCategoryID = c.ID
		} else {
			form.CategoryIDs = append(form.CategoryIDs, c.ID)
		}
	}

	pageData := models.PageData{
		User: user,
		PageContent: PostFormPage{
			PostID:     id,
			Form:       form,
			Categories: cats,
		},
	}
	utils.RenderTemplate(w, http.StatusOK, "post_form", pageData)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user_session").(*models.UserInfo)
	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		handleError(w, customerrors.ErrBadRequest)
		return
	}

	form, ok := h.parsePostForm(w, r)
	if !ok {
		return
	}

	if !form.Validate() {
		cats, err := h.catRep.GetAllCategories(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}
		pageData := models.PageData{
			User: user,
			PageContent: PostFormPage{
				PostID: id,
				Form: PostForm{
					Title:   form.Title,
					Content: form.Content,
				},
				Categories: cats,
			},
		}
		utils.RenderTemplate(w, http.StatusOK, "post_form", pageData)
		return
	}

	postUpdate := models.PostUpdate{
		ID:             id,
		Title:          form.Title,
		Content:        form.Content,
		MainCategoryID: form.MainCategoryID,
		CategoryIDs:    form.CategoryIDs,
	}

	err = h.postRep.UpdatePost(r.Context(), postUpdate, user.Id)
	if err != nil {
		handleError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)

}
