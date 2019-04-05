package controllers

import (
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"

	"github.com/hackallcode/db-homework/internal/pkg/db"
	"github.com/hackallcode/db-homework/internal/pkg/models"
)

func ClearHandler(ctx *routing.Context) error {
	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	err = db.TruncateAll(tx)
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	err = tx.Commit()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	SendJson(ctx, fasthttp.StatusOK, models.NewMessageAnswer("All data was successfully deleted!"))
	return nil
}

func StatusHandler(ctx *routing.Context) error {
	tx, err := db.Begin()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}
	defer tx.Rollback()

	statusData, err := db.Status(tx)
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	err = tx.Commit()
	if err != nil {
		SendJson(ctx, fasthttp.StatusInternalServerError, models.NewErrorAnswer(err.Error()))
		return nil
	}

	SendJson(ctx, fasthttp.StatusOK, statusData)
	return nil
}
