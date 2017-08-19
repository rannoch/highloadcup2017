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
	var fromDate = ctx.QueryArgs().GetUintOrZero("fromDate")
	var toDate = ctx.QueryArgs().GetUintOrZero("toDate")
	var country = (string)(ctx.QueryArgs().Peek("country"))
	var toDistance = ctx.QueryArgs().GetUintOrZero("toDistance")

	var visits []models.Visit
	var conditions []storage.Condition
	var joins []storage.Join

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		log.Printf("id parse error %v \n", ctx.UserValue("id"))
		return
	}

	idCondition := storage.Condition{
		Param:         "user",
		Value:         strconv.Itoa(id),
		Operator:      "=",
		JoinCondition: "and",
	}
	conditions = append(conditions, idCondition)

	if fromDate > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "visited_at ",
			Value:         strconv.Itoa(fromDate),
			Operator:      ">",
			JoinCondition: "and",
		})
	}

	if toDate > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "visited_at ",
			Value:         strconv.Itoa(toDate),
			Operator:      "<",
			JoinCondition: "and",
		})
	}

	if toDistance > 0 || len(country) > 0 {
		joins = append(joins, storage.Join{
			Name: "location",
			Type: "left",
			Condition: storage.Condition{
				Param:    "visit.location",
				Value:    "location.id",
				Operator: "=",
			},
		})
	}

	if len(country) > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "country",
			Value:         "'" +country + "'",
			Operator:      "=",
			JoinCondition: "and",
		})
	}

	if toDistance > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "distance",
			Value:         strconv.Itoa(toDistance),
			Operator:      "<",
			JoinCondition: "and",
		})
	}

	err = storage.Db.SelectEntityMultiple(&visits, (&models.Visit{}).GetFields("visit"), joins, conditions)

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
