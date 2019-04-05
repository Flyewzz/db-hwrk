package controllers

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func CreateThreadHandler(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	threadData := models.Thread{}
	err := threadData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	threadData.Forum = slug

	if !threadData.Validate() {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	err = db.CreateThread(tx, &threadData)
	if err != nil {
		if err == models.AlreadyExistsError {
			SendJson(ctx, fasthttp.StatusConflict, threadData)
		} else if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.ForumNotFoundAnswer)
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

	SendJson(ctx, fasthttp.StatusCreated, threadData)
	return nil
}

func GetThreadHandler(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	if slugOrId == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	forumData, err := db.GetThread(tx, slugOrId)
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

func UpdateThreadHandler(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	if slugOrId == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	updateData := models.ThreadUpdate{}
	err := updateData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	updateData.SlugOrId = slugOrId

	if !updateData.Validate() {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	threadData, err := db.UpdateThread(tx, &updateData)
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

	SendJson(ctx, fasthttp.StatusOK, threadData)
	return nil
}

func VoteThreadHandler(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	if slugOrId == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	voteData := models.Vote{}
	err := voteData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	voteData.SlugOrId = slugOrId

	if !voteData.Validate() {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	threadData, err := db.VoteThread(tx, &voteData)
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

	SendJson(ctx, fasthttp.StatusOK, threadData)
	return nil
}

func GetForumThreadsHandler(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	threads, err := db.GetForumThreads(tx, slug, limit, since, desc)
	if err != nil {
		if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.ForumNotFoundAnswer)
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

	SendJson(ctx, fasthttp.StatusOK, threads)
	return nil
}
