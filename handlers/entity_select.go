package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/storage"
	"github.com/rannoch/highloadcup2017/models"
	"strings"
	"database/sql"
	"encoding/json"
	"log"
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
		log.Printf("id parse error %v \n", ctx.UserValue("id"))
		return
	}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
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
		Param:    "id",
		Value:    strconv.Itoa(id),
		Operator: "=",
		Join:     "and",
	}
	conditions = append(conditions, idCondition)

	err = storage.Db.SelectEntity(entity, conditions)

	if err == sql.ErrNoRows {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	response, err := json.Marshal(entity)
	if err != nil {
		log.Println(err)
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
