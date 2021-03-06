package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/mysql_server/storage"
	"github.com/rannoch/highloadcup2017/mysql_server/models"
	"strings"
	"database/sql"
	"encoding/json"
	"strconv"
)

func EntitySelectHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var entityValue string
	var id int
	var entity storage.Entity
	var conditions []storage.Condition

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusNotFound)
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
	ctx.UserValue("entity")

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

	response, err := json.Marshal(entity)
	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
