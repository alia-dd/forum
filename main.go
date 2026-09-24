package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gitea.kood.tech/jyrkikarhunen/forum/database"
	"gitea.kood.tech/jyrkikarhunen/forum/handlers"
	"gitea.kood.tech/jyrkikarhunen/forum/middleware"
	"gitea.kood.tech/jyrkikarhunen/forum/models"
	"gitea.kood.tech/jyrkikarhunen/forum/repository"
	"gitea.kood.tech/jyrkikarhunen/forum/service"
	"gitea.kood.tech/jyrkikarhunen/forum/utils"
)

func main() {
	var deleteDB bool
	var seedDB bool

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-d" || os.Args[i] == "-delete" {
			deleteDB = true
		}
		if os.Args[i] == "-i" || os.Args[i] == "-init" {
			seedDB = true
		}
	}

	if deleteDB {
		err := os.Remove("data/forum.db")
		if err != nil {
			fmt.Println(err)
		}
	}

	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if seedDB {
		err := database.SeedData(db)
		if err != nil {
			fmt.Println(err)
		}
	}

	utils.InitializeTemplate()

	mux := http.NewServeMux()

	sessionRepo := repository.NewSessionRepository(db)
	userRepo := repository.NewUserRepository(db)

	reactRepo := repository.NewReactionRepositry(db)

	postRepo := repository.NewPostRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	commentRepo := repository.NewCommentRepository(db)

	userService := service.NewUserService(userRepo, sessionRepo)
	service.NewReactionHandler(reactRepo)

	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postRepo, categoryRepo, commentRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)

	searchRepo := repository.NewSearchRepository(db)
	searchHandler := handlers.NewSearchHandler(searchRepo)

	mux.Handle("GET /static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// route check if the routes bellow are not called and difult to this one and if not calls 404 page
	// wrappes the whole with call with the safeguarded other while the user session will not we recognize and will log you out
	mux.HandleFunc("GET /", middleware.Recoverer(
		middleware.AllowGuest(sessionRepo, userRepo, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				user, ok := r.Context().Value("user_session").(*models.UserInfo)
				if !ok {
					http.Redirect(w, r, "/user/login", http.StatusSeeOther)
					return
				}

				pageData := models.PageData[models.ErrorStruct]{
					User:        user,
					IsOwner:     ok,
					PageContent: models.ErrorStruct{Error: "404", ErrorMs: "Page Not Found"},
				}
				utils.RenderTemplate(w, http.StatusNotFound, "error", pageData)
				return
			}

			(postHandler.GetPosts)(w, r)

		}),
	))

	// Profile page is user specific and is safeguarded by the Restrict middleware
	// if there is no active session, it redirects to the login page.
	mux.HandleFunc("GET /user/profile", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, handlers.Profile)))
	mux.HandleFunc("GET /user/profile/{username}", middleware.Recoverer(middleware.AllowGuest(sessionRepo, userRepo, userHandler.GetOtherUserProfile)))

	mux.HandleFunc("GET /user/profile/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, userHandler.GetEditUserProfile)))
	mux.HandleFunc("POST /user/profile/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, userHandler.UpdateUserProfile)))

	mux.HandleFunc("GET /user/register", middleware.Recoverer(userHandler.GetRegisterUser))
	mux.HandleFunc("POST /user/register", middleware.Recoverer(userHandler.PostRegisterUser))

	mux.HandleFunc("GET /user/login", middleware.Recoverer(userHandler.GetSignInUser))
	mux.HandleFunc("POST /user/login", middleware.Recoverer(userHandler.SignInUser))

	mux.HandleFunc("GET /user/profile/password", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, userHandler.GetChangePassword)))
	mux.HandleFunc("POST /user/profile/password", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, userHandler.PostChangePassword)))

	mux.HandleFunc("POST /user/logout", middleware.Recoverer(userHandler.SignOutUser))

	// mux.HandleFunc("GET /api/user/check-username", middleware.Recoverer(userHandler.CheckIfAvailabe))
	// mux.HandleFunc("GET /api/user/check-email", middleware.Recoverer(userHandler.CheckIfAvailabe))

	// GET /post/{id} - just an int
	mux.HandleFunc("GET /post/{id}", middleware.Recoverer(middleware.AllowGuest(sessionRepo, userRepo, postHandler.GetPostByID)))

	//below need auth
	// GET /post/new - loads template (once implemented) for submitting post
	mux.HandleFunc("GET /post/new", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, postHandler.NewPostForm)))
	// POST /post/new - create post. Form: title, content, main_category={id}, category={id}… (sides)
	mux.HandleFunc("POST /post/new", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, postHandler.CreatePost)))
	// GET /post/{id}/edit - loads template (once implemented) for editing a post you own
	mux.HandleFunc("GET /post/{id}/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, postHandler.EditPostForm)))
	// POST /post/{id}/edit - update your post. Form: title, content, main_category={id}, category={id}… (sides)
	mux.HandleFunc("POST /post/{id}/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, postHandler.UpdatePost)))
	mux.HandleFunc("POST /post/{id}/delete", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, postHandler.DeletePost)))

	// GET /category/new - loads template (once implemented) for creating new category
	mux.HandleFunc("GET /category/new", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, categoryHandler.NewCategoryForm)))
	// POST /category/new - create category. Form: name
	mux.HandleFunc("POST /category/new", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, categoryHandler.CreateCategory)))
	// GET /category/{id}/edit - loads template (once implemented) for editing existing category
	mux.HandleFunc("GET /category/{id}/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, categoryHandler.EditCategoryForm)))
	// POST /category/{id}/edit - rename category. Form: name
	mux.HandleFunc("POST /category/{id}/edit", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, categoryHandler.UpdateCategory)))
	// POST /category/{id}/delete - delete category (if it has no references elsewhere)
	mux.HandleFunc("POST /category/{id}/delete", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, categoryHandler.DeleteCategory)))

	// POST /post/{id}/comment/create - for commenting
	mux.HandleFunc("POST /post/{id}/comment/create", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, commentHandler.CreateComment)))
	// GET /post/{id} - for commenting
	// mux.HandleFunc("GET /post/{id}/comment", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, commentHandler.EditComment)))
	// GET /post/{id} - for commenting
	// mux.HandleFunc("GET /post/{id}/comment", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, commentHandler.UpdateComment)))
	// POST /post/{id} - for commenting
	// mux.HandleFunc("POST /post/{id}/comment/delete", middleware.Recoverer(middleware.Restrict(sessionRepo, userRepo, commentHandler.DeleteComment)))
	//above need auth

	mux.HandleFunc("GET /search", middleware.Recoverer(middleware.AllowGuest(sessionRepo, userRepo, searchHandler.Search)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("Server starting on http://localhost:8080")

	log.Fatal(server.ListenAndServe())
}

// self note
// delete session after logout
// is the session relly working needs more test
