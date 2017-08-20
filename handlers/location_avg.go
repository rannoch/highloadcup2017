package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/storage"
	"github.com/rannoch/highloadcup2017/models"
	"strconv"
	"database/sql"
	"encoding/json"
	"time"
)

func LocationsAvgHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id, fromDate, toDate, fromAge, toAge int
	var err error

	if ctx.QueryArgs().Has("fromDate") {
		fromDate, err = ctx.QueryArgs().GetUint("fromDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}
	if ctx.QueryArgs().Has("toDate") {
		toDate, err = ctx.QueryArgs().GetUint("toDate")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}
	if ctx.QueryArgs().Has("fromAge") {
		fromAge, err = ctx.QueryArgs().GetUint("fromAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}
	if ctx.QueryArgs().Has("toAge") {
		toAge, err = ctx.QueryArgs().GetUint("toAge")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var gender = (string)(ctx.QueryArgs().Peek("gender"))

	var avg float32 = 0
	var conditions []storage.Condition
	var joins []storage.Join

	id, err = strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	if gender != "" && !(gender == "m" || gender == "f") {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	idCondition := storage.Condition{
		Param:         "visit.location",
		Value:         strconv.Itoa(id),
		Operator:      "=",
		JoinCondition: "and",
	}
	conditions = append(conditions, idCondition)

	location := models.Location{}
	err = storage.Db.SelectEntity(&location, []storage.Condition{
		{
			Param:         "id",
			Value:         strconv.Itoa(id),
			Operator:      "=",
			JoinCondition: "and",
		},
	})

	if err == sql.ErrNoRows {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

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
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}