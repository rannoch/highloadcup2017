package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/antonholmquist/jason"
	"encoding/json"
	"sort"
	"bytes"
	"github.com/rannoch/highloadcup2017/memory/models"
	"github.com/rannoch/highloadcup2017/memory/storage"
)

func EntitityNewHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var entityValue string
	var params map[string]interface{}

	defer ctx.SetConnectionClose()

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch entityValue{
	case "users":
		entity := &models.User{}
		// check params
		postBody, err := jason.NewValueFromBytes(ctx.PostBody())

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		params, ok := postBody.Interface().(map[string]interface{})
		if !ok {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		if !entity.ValidateParams(params, "insert") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)

		storage.UserDb[entity.Id] = entity
		storage.UserCount++
	case "locations":
		entity := &models.Location{}

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

		storage.LocationDb[entity.Id] = entity
		storage.LocationCount++
	case "visits":
		entity := &models.Visit{}

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
		storage.VisitDb[entity.Id] = entity
		storage.VisitCount++

		user := storage.UserDb[entity.User]
		location := storage.LocationDb[entity.Location]

		entity.User_model = user
		entity.Location_model = location

		user.Visits = append(user.Visits, entity)
		location.Visits = append(location.Visits, entity)

		sort.Sort(models.VisitByDateAsc(user.Visits))
	}

	buffer := bufPool.Get().(*bytes.Buffer)
	buffer.Reset()
	buffer.WriteString("{}")

	ctx.Write(buffer.Bytes())
	bufPool.Put(buffer)
}

