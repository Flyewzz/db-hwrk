package controllers

import (
	"log"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func SendJson(ctx *routing.Context, statusCode int, model models.OutputModel) {
	json, err := model.MarshalJSON()
	if err != nil {
		log.Println(err)
		return
	}

	ctx.SetContentType("application/json; charset=utf-8")
	ctx.Response.Header.SetStatusCode(statusCode)
	ctx.SetBody(json)
}

func ApiHandler(ctx *routing.Context) error {
	SendJson(ctx, fasthttp.StatusOK, models.NewMessageAnswer("It's api for forum!"))
	return nil
}
