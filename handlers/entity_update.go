package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/storage"
	"strconv"
	"strings"
	"github.com/rannoch/highloadcup2017/models"
	"encoding/json"
	"database/sql"
)

func EntityUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int
	var entity storage.Entity
	var conditions []storage.Condition
	var entityValue string
	var params map[string]interface{}

	if ctx.UserValue("id").(string) == "new" {
		EntitityNewHandler(ctx)
		return
	}

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch {
	case strings.Contains(entityValue, "user"):
		entity = &models.User{}
	case strings.Contains(entityValue, "location"):
		entity = &models.Location{}
	case strings.Contains(entityValue, "visit"):
		entity = &models.Visit{}
	}

	// check params
	err = json.Unmarshal(ctx.PostBody(), &params)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if !entity.ValidateParams(params, "update") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	idCondition := storage.Condition{
		Param:         "id",
		Value:         strconv.Itoa(id),
		Operator:      "=",
		JoinCondition: "and",
	}
	conditions = append(conditions, idCondition)

	err = storage.Db.SelectEntity(entity, conditions)

	if err == sql.ErrNoRows {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	_, err = storage.Db.UpdateEntity(entity, params, conditions)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	ctx.SetBody([]byte("{}"))
}

func EntitityNewHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var entityValue string
	var entity storage.Entity
	var params map[string]interface{}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch {
	case strings.Contains(entityValue, "user"):
		entity = &models.User{}
	case strings.Contains(entityValue, "location"):
		entity = &models.Location{}
	case strings.Contains(entityValue, "visit"):
		entity = &models.Visit{}
	}

	// check params
	err := json.Unmarshal(ctx.PostBody(), &params)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	if !entity.ValidateParams(params, "insert") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	entity.SetParams(params)

	err = storage.Db.InsertEntity(entity)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	ctx.SetBody([]byte("{}"))
}