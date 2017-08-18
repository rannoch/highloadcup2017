package handlers

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"log"
	"github.com/rannoch/highloadcup2017/storage"
	"database/sql"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/models"
)

func UsersVisitsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int
	/*var fromDate *int
	var toDate *int
	var country *string
	var toDistance *int*/
	var visits []models.Visit
	//var entities []storage.Entity
	var conditions []storage.Condition

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		log.Printf("id parse error %v \n", ctx.UserValue("id"))
		return
	}

	idCondition := storage.Condition{
		Param:    "user",
		Value:    strconv.Itoa(id),
		Operator: "=",
		Join:     "and",
	}
	conditions = append(conditions, idCondition)

	err = storage.Db.SelectEntityMultiple(&visits, conditions)

	if err == sql.ErrNoRows {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	response, err := json.Marshal(visits)
	if err != nil {
		log.Println(err)
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}