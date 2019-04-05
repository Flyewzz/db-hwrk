package controllers

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func CreateUserHandler(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	if nickname == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	userData := models.User{}
	err := userData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	userData.Nickname = nickname

	if !userData.Validate() {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectDataAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	usersData, err := db.CreateUser(tx, userData)
	if err != nil {
		if err == models.AlreadyExistsError {
			SendJson(ctx, fasthttp.StatusConflict, usersData)
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

	SendJson(ctx, fasthttp.StatusCreated, userData)
	return nil
}

func GetUserHandler(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	if nickname == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	userData, err := db.GetUser(tx, nickname)
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

	SendJson(ctx, fasthttp.StatusOK, userData)
	return nil
}

func UpdateUserHandler(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	if nickname == "" {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectUrlAnswer)
		return nil
	}

	updateData := models.UserUpdate{}
	err := updateData.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		SendJson(ctx, fasthttp.StatusBadRequest, models.IncorrectJsonAnswer)
		return nil
	}
	updateData.Nickname = nickname

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

	err = db.UpdateUser(tx, &updateData)
	if err != nil {
		if err == models.NotFoundError {
			SendJson(ctx, fasthttp.StatusNotFound, models.UserNotFoundAnswer)
		} else if err == models.AlreadyExistsError {
			SendJson(ctx, fasthttp.StatusConflict, models.DuplicateDataAnswer)
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

	SendJson(ctx, fasthttp.StatusOK, models.User{
		Nickname: updateData.Nickname,
		Email:    updateData.Email,
		FullName: updateData.FullName,
		About:    updateData.About,
	})
	return nil
}

func GetForumUsersHandler(ctx *routing.Context) error {
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

	threads, err := db.GetForumUsers(tx, slug, limit, since, desc)
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
