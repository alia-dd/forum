package server

import (
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func registerRoutes(mux *http.ServeMux, d routeDependencies) {
	mux.Handle("GET /static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// Wraps the route with guest/session handling and renders a 404 page
	// for any path other than "/". If the session is invalid, the user
	// may be redirected or logged out, depending on handleGuest's behavior.
	handleGuest(mux, "GET /", d, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			user, _ := r.Context().Value("user_session").(*models.UserInfo)
			pageData := models.PageData[models.ErrorStruct]{
				User: user,
				PageContent: models.ErrorStruct{
					Error:   "404",
					ErrorMs: "Page Not Found",
				},
			}
			utils.RenderTemplate(w, http.StatusNotFound, "error", pageData)
			return
		}
		d.postHandler.GetPosts(w, r)
	})

	// Admin
	handleAdmin(mux, "GET /admin/confirmation", d, d.userHandler.GetConfirmUserAction)
	handleAdmin(mux, "GET /admin/users", d, d.userHandler.GetAdminUserTable)
	handleAdmin(mux, "GET /admin/posts", d, d.userHandler.GetAdminPostTable)
	handleAdmin(mux, "GET /admin/categories", d, d.userHandler.GetCategoryTable)

	handleAdmin(mux, "POST /admin/users/action", d, d.userHandler.PostAdminUserAction)
	handleAdmin(mux, "POST /admin/post/action", d, d.userHandler.PostAdminPostrAction)

	// User profiles
	handleRestricted(mux, "GET /user/profile", d, d.userHandler.GetOwnProfile)
	handleGuest(mux, "GET /user/profile/{username}", d, d.userHandler.GetOtherUserProfile)

	handleRestricted(mux, "GET /user/profile/edit", d, d.userHandler.GetEditUserProfile)
	handleRestricted(mux, "POST /user/profile/edit", d, d.userHandler.UpdateUserProfile)

	handleRestricted(mux, "GET /user/profile/password", d, d.userHandler.GetChangePassword)
	handleRestricted(mux, "POST /user/profile/password", d, d.userHandler.PostChangePassword)

	handleRestricted(mux, "POST /user/profile/delete", d, d.userHandler.DeleteUserAcccount)

	// Auth
	handle(mux, "GET /user/register", d.userHandler.GetRegisterUser)
	handle(mux, "POST /user/register", d.userHandler.PostRegisterUser)

	handle(mux, "GET /user/login", d.userHandler.GetSignInUser)
	handle(mux, "POST /user/login", d.userHandler.SignInUser)
	handle(mux, "POST /user/logout", d.userHandler.SignOutUser)

	// Reactions
	handleRestricted(mux, "POST /react", d, d.reactHandler.Reaction)

	// Posts
	handleGuest(mux, "GET /post/{id}", d, d.postHandler.GetPostByID)

	handleRestricted(mux, "GET /post/new", d, d.postHandler.NewPostForm)
	handleRestricted(mux, "POST /post/new", d, d.postHandler.CreatePost)
	handleRestricted(mux, "GET /post/{id}/edit", d, d.postHandler.EditPostForm)
	handleRestricted(mux, "POST /post/{id}/edit", d, d.postHandler.UpdatePost)
	handleRestricted(mux, "POST /post/{id}/delete", d, d.postHandler.DeletePost)

	// Categories
	handleAdmin(mux, "GET /category/new", d, d.categoryHandler.NewCategoryForm)
	handleAdmin(mux, "POST /category/new", d, d.categoryHandler.CreateCategory)
	handleAdmin(mux, "GET /category/{id}/edit", d, d.categoryHandler.EditCategoryForm)
	handleAdmin(mux, "POST /category/{id}/edit", d, d.categoryHandler.UpdateCategory)
	handleAdmin(mux, "POST /category/{id}/delete", d, d.categoryHandler.DeleteCategory)

	// Comments
	handleRestricted(mux, "POST /post/{id}/comment/create", d, d.commentHandler.CreateComment)
	handleRestricted(mux, "GET /post/{postID}/comment/{commentID}/edit", d, d.commentHandler.GetCommentEditForm)
	handleRestricted(mux, "PUT /post/{postID}/comment/{commentID}", d, d.commentHandler.UpdateComment)
	handleRestricted(mux, "GET /post/{postID}/comment/{commentID}", d, d.commentHandler.GetComment)
	handleRestricted(mux, "DELETE /post/{postID}/comment/{commentID}/delete", d, d.commentHandler.DeleteComment)

	handleRestricted(mux, "GET /post/{postID}/comment/{commentID}/reply", d, d.commentHandler.GetReplyForm)
	handleRestricted(mux, "POST /post/{postID}/comment/{commentID}/reply", d, d.commentHandler.CreateReply)

	handleGuest(mux, "GET /post/{postID}/comment/{commentID}/replies", d, d.commentHandler.GetCommentReplies)

	// Search
	handleGuest(mux, "GET /search", d, d.searchHandler.Search)
}
