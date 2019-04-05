package controllers

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func CreateForumHandler(ctx *routing.Context) error {
	forumData := models.Forum{}
	err := forumData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}

	if !forumData.Validate() {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	err = db.CreateForum(tx, &forumData)
	if err != nil {
		if err == models.AlreadyExistsError {
			SendJson(ctx, fasthttp.StatusConflict, forumData)
		} else if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.UserNotFoundAnswer)
		} else {
			SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		}
		return nil
	}

	err = tx.Commit()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	SendJson(ctx, fasthttp.StatusCreated, forumData)
	return nil
}

func GetForumHandler(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	forumData, err := db.GetForum(tx, slug)
	if err != nil {
		if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.UserNotFoundAnswer)
		} else {
			SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		}
		return nil
	}

	err = tx.Commit()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	SendJson(ctx, fasthttp.StatusOK, forumData)
	return nil
}
