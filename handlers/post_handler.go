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
	user := r.Context().Value("user_session").(*models.UserInfo)

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

	fmt.Fprintf(w, "posts=%v cats=%v", post, cats)

	//execute postviewtemplate with user, post
}

func (h *PostHandler) NewPostForm(w http.ResponseWriter, r *http.Request) {
	cats, err := h.catRep.GetAllCategories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	fmt.Fprintf(w, "cats=%v", cats)

	//execute createposttemplate with user, cats
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user_session").(*models.UserInfo)
	form, ok := h.parsePostForm(w, r)
	if !ok {
		return
	}

	if !form.Validate() {
		//execute same template with form values (prefilled) and form errors next to relevant sections
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
	user := r.Context().Value("user_session").(*models.UserInfo)

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

	fmt.Fprintf(w, "post=%v cats=%v", post, cats)

	//execute postedittemplate with user, post, cats
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user_session").(*models.UserInfo)

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
		//execute same template with form values (prefilled) and form errors next to relevant sections
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
