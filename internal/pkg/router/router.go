package router

import (
	"log"
	"strings"

	"github.com/qiangxue/fasthttp-routing"

	c "github.com/hackallcode/db-homework/internal/pkg/controllers"
)

type Route struct {
	Path    string
	Method  string
	Handler routing.Handler
}

type Group struct {
	Prefix string
	Routes []Route
	Groups []Group
}

var routes = Group{
	Routes: []Route{
		{Path: "/", Method: strings.Join(routing.Methods, ","), Handler: c.ApiHandler},
	},
	Groups: []Group{
		{Prefix: "/user", Routes: []Route{
			{Path: "/<nickname>/create", Method: "POST", Handler: c.CreateUserHandler},
			{Path: "/<nickname>/profile", Method: "GET", Handler: c.GetUserHandler},
			{Path: "/<nickname>/profile", Method: "POST", Handler: c.UpdateUserHandler},
		}},
		{Prefix: "/forum", Routes: []Route{
			{Path: "/create", Method: "POST", Handler: c.CreateForumHandler},
			{Path: "/<slug>/details", Method: "GET", Handler: c.GetForumHandler},
			{Path: "/<slug>/create", Method: "POST", Handler: c.CreateThreadHandler},
			{Path: "/<slug>/threads", Method: "GET", Handler: c.GetForumThreadsHandler},
			{Path: "/<slug>/users", Method: "GET", Handler: c.GetForumUsersHandler},
		}},
		{Prefix: "/post", Routes: []Route{
			{Path: "/<id>/details", Method: "GET", Handler: c.GetFullPostHandler},
			{Path: "/<id>/details", Method: "POST", Handler: c.UpdatePostHandler},
		}},
		{Prefix: "/thread", Routes: []Route{
			{Path: "/<slug_or_id>/create", Method: "POST", Handler: c.CreatePostsHandler},
			{Path: "/<slug_or_id>/posts", Method: "GET", Handler: c.GetThreadPostsHandler},
			{Path: "/<slug_or_id>/details", Method: "GET", Handler: c.GetThreadHandler},
			{Path: "/<slug_or_id>/details", Method: "POST", Handler: c.UpdateThreadHandler},
			{Path: "/<slug_or_id>/vote", Method: "POST", Handler: c.VoteThreadHandler},
		}},
		{Prefix: "/service", Routes: []Route{
			{Path: "/clear", Method: "POST", Handler: c.ClearHandler},
			{Path: "/status", Method: "GET", Handler: c.StatusHandler},
		}},
	},
}

func initGroup(router *routing.RouteGroup, group Group) {
	routeGroup := router.Group(group.Prefix)
	for _, route := range group.Routes {
		routeGroup.To(strings.ToUpper(route.Method), route.Path, route.Handler)
	}
	for _, child := range group.Groups {
		initGroup(routeGroup, child)
	}
}

func InitRouter(prefix string) *routing.Router {
	routes.Prefix = prefix

	router := routing.New()
	initGroup(&router.RouteGroup, routes)
	log.Println("router has been initialized")

	return router
}
