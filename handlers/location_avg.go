package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/storage"
	"github.com/rannoch/highloadcup2017/models"
	"log"
	"strconv"
	"database/sql"
	"encoding/json"
	"time"
)

func LocationsAvgHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int

	var fromDate = ctx.QueryArgs().GetUintOrZero("fromDate")
	var toDate = ctx.QueryArgs().GetUintOrZero("toDate")
	var fromAge = ctx.QueryArgs().GetUintOrZero("fromAge")
	var toAge = ctx.QueryArgs().GetUintOrZero("toAge")
	var gender = (string)(ctx.QueryArgs().Peek("gender"))

	var avg float32 = 0
	var conditions []storage.Condition
	var joins []storage.Join

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		log.Printf("id parse error %v \n", ctx.UserValue("id"))
		return
	}

	if gender != "" && !(gender == "m" || gender == "f") {
		ctx.Error("", fasthttp.StatusBadRequest)
		log.Printf("invalid gender %s \n", gender)
		return
	}

	idCondition := storage.Condition{
		Param:         "visit.location",
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

	if fromAge > 0 || toAge > 0 || len(gender) > 0 {
		joins = append(joins, storage.Join{
			Name: "user",
			Type: "left",
			Condition: storage.Condition{
				Param:    "visit.user",
				Value:    "user.id",
				Operator: "=",
			},
		})
	}

	if fromAge > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "user.birth_date ",
			Value:         strconv.Itoa(int(time.Now().AddDate(-fromAge, 0, 0).Unix())),
			Operator:      "<",
			JoinCondition: "and",
		})
	}

	if toAge > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "user.birth_date ",
			Value:         strconv.Itoa(int(time.Now().AddDate(-toAge, 0, 0).Unix())),
			Operator:      ">",
			JoinCondition: "and",
		})
	}

	if len(gender) > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "user.gender ",
			Value:         "'" + gender + "'",
			Operator:      "=",
			JoinCondition: "and",
		})
	}

	avg, err = storage.Db.GetAverage(&models.Visit{}, "mark", joins, conditions)

	if err == sql.ErrNoRows {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	response, err := json.Marshal(map[string]interface{}{"avg" : models.FloatPrecision5(avg) })
	if err != nil {
		log.Println(err)
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}