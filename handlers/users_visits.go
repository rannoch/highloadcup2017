package handlers

import (
	"github.com/valyala/fasthttp"
	"strconv"
	"github.com/rannoch/highloadcup2017/storage"
	"database/sql"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/models"
)

func UsersVisitsHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int
	var fromDate, toDate, toDistance int
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

	if ctx.QueryArgs().Has("toDistance") {
		toDistance, err = ctx.QueryArgs().GetUint("toDistance")

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}
	}

	var country = (string)(ctx.QueryArgs().Peek("country"))

	var visits []models.Visit = []models.Visit{}
	var conditions []storage.Condition
	var joins []storage.Join

	id, err = strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusNotFound)
		return
	}

	idCondition := storage.Condition{
		Param:         "user",
		Value:         strconv.Itoa(id),
		Operator:      "=",
		JoinCondition: "and",
	}
	conditions = append(conditions, idCondition)

	user := models.User{}
	err = storage.Db.SelectEntity(&user, []storage.Condition{
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

	joins = append(joins, storage.Join{
		Name: "location",
		Type: "left",
		Condition: storage.Condition{
			Param:    "visit.location",
			Value:    "location.id",
			Operator: "=",
		},
	})

	if len(country) > 0 {
		conditions = append(conditions, storage.Condition{
			Param:         "country",
			Value:         "'" + country + "'",
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

	err = storage.Db.SelectEntityMultiple(&visits, []string{}, joins, conditions, storage.Sort{Fields:[]string{"visited_at"}, Direction:"asc"})

	if err == sql.ErrNoRows {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	visitsResponse := []interface{}{}

	for _, visit := range visits {
		v := map[string]interface{}{
			"mark":       visit.Mark,
			"visited_at": visit.Visited_at,
			"place":      visit.LocationChild.Place,
		}

		visitsResponse = append(visitsResponse, v)
	}

	response, err := json.Marshal(map[string]interface{}{"visits": visitsResponse})
	if err != nil {
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		return
	}

	ctx.SetBody(response)
}
