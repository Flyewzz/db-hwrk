package server

import (
	"log"
	"strconv"

	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/router"
)

type Params struct {
	Port  int64
	Url   string
	Reset bool
}

func StartApp(params Params) error {
	portStr := strconv.FormatInt(params.Port, 10)

	if err := db.Open(params.Reset); err != nil {
		return err
	}

	apiRouter := router.InitRouter(params.Url)

	log.Printf("server is starting up at http://localhost:%v%v/...\n", portStr, params.Url)
	return fasthttp.ListenAndServe(":"+portStr, apiRouter.HandleRequest)
}

func StopApp() error {
	log.Println("server is stopping...")
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
