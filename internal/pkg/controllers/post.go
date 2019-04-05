package controllers

import (
	"strconv"
	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func CreatePostsHandler(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	if slugOrId == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	postsData := models.Posts{}
	err := postsData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}

	// if !postsData.Validate() {
	// 	SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
	// 	return nil
	// }

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	err = db.CreatePosts(tx, slugOrId, &postsData)
	if err != nil {
		if err == models.AlreadyExistsError {
			SendJson(ctx, fasthttp.StatusConflict, postsData)
		} else if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.ForumNotFoundAnswer)
		} else if err == models.ConflictDataError {
			SendJson(ctx, fasthttp.StatusConflict, models.ParentIncorrectAnswer)
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

	SendJson(ctx, fasthttp.StatusCreated, postsData)
	return nil
}

func GetThreadPostsHandler(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	if slugOrId == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	limit := string(ctx.QueryArgs().Peek("limit"))
	since := string(ctx.QueryArgs().Peek("since"))
	sort := string(ctx.QueryArgs().Peek("sort"))
	desc := string(ctx.QueryArgs().Peek("desc"))

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	threads, err := db.GetThreadPosts(tx, slugOrId, limit, since, sort, desc)
	if err != nil {
		if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.ThreadNotFoundAnswer)
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

func UpdatePostHandler(ctx *routing.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if id == 0 || err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	updateData := models.PostUpdate{}
	err = updateData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	updateData.Id = id

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

	postData, err := db.UpdatePost(tx, &updateData)
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

	SendJson(ctx, fasthttp.StatusOK, postData)
	return nil
}

func GetFullPostHandler(ctx *routing.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if id == 0 || err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	related := strings.Split(string(ctx.QueryArgs().Peek("related")), ",")

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	postData, err := db.GetFullPost(tx, id, related)
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

	SendJson(ctx, fasthttp.StatusOK, postData)
	return nil
}
