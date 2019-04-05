package server

import (
	"log"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	c "github.com/hackallcode/db-homework/internal/pkg/controllers"
	"github.com/hackallcode/db-homework/internal/pkg/db"
)

type Params struct {
	Port      string
	ApiPrefix string
}

func StartApp(params Params) error {
	log.Println("Server starting at " + params.Port)

	if err := db.Open(); err != nil {
		return err
	}

	router := routing.New()

	apiRouter := router.Group(params.ApiPrefix)

	apiRouter.Get("/", c.ApiHandler)

	userRouter := apiRouter.Group("/user")
	userRouter.Post("/<nickname>/create", c.CreateUserHandler)
	userRouter.Get("/<nickname>/profile", c.GetUserHandler)
	userRouter.Post("/<nickname>/profile", c.UpdateUserHandler)

	forumRouter := apiRouter.Group("/forum")
	forumRouter.Post("/create", c.CreateForumHandler)
	forumRouter.Get("/<slug>/details", c.GetForumHandler)
	forumRouter.Post("/<slug>/create", c.CreateThreadHandler)
	forumRouter.Get("/<slug>/threads", c.GetForumThreadsHandler)
	forumRouter.Get("/<slug>/users", c.GetForumUsersHandler)

	postRouter := apiRouter.Group("/post")
	postRouter.Get("/<id>/details", c.GetFullPostHandler)
	postRouter.Post("/<id>/details", c.UpdatePostHandler)

	threadRouter := apiRouter.Group("/thread")
	threadRouter.Post("/<slug_or_id>/create", c.CreatePostsHandler)
	threadRouter.Get("/<slug_or_id>/posts", c.GetThreadPostsHandler)
	threadRouter.Get("/<slug_or_id>/details", c.GetThreadHandler)
	threadRouter.Post("/<slug_or_id>/details", c.UpdateThreadHandler)
	threadRouter.Post("/<slug_or_id>/vote", c.VoteThreadHandler)

	serviceRouter := apiRouter.Group("/service")
	serviceRouter.Post("/clear", c.ClearHandler)
	serviceRouter.Get("/status", c.StatusHandler)

	return fasthttp.ListenAndServe(":"+params.Port, router.HandleRequest)
}

func StopApp() error {
	log.Println("Stopping server...")
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
