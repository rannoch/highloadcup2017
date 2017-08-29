package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/antonholmquist/jason"
	"encoding/json"
	"sort"
	"github.com/rannoch/highloadcup2017/epoll_server/models"
	"github.com/rannoch/highloadcup2017/epoll_server/storage"
	"github.com/rannoch/highloadcup2017/epoll_server/server"
)

func EntitityNewHandler(ctx *server.HlcupCtx, entityValue string) {
	var params map[string]interface{}

	switch entityValue{
	case "users":
		entity := &models.User{}
		// check params
		postBody, err := jason.NewValueFromBytes(ctx.PostBody)

		if err != nil {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		params, ok := postBody.Interface().(map[string]interface{})
		if !ok {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		if !entity.ValidateParams(params, "insert") {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)

		storage.UserDb[entity.Id] = entity
		storage.UserCount++
		//storage.UserBytesDb[entity.Id] = entity.GetBytes()
	case "locations":
		entity := &models.Location{}

		err := json.Unmarshal(ctx.PostBody, &params)
		if err != nil {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		if !entity.ValidateParams(params, "insert") {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)

		storage.LocationDb[entity.Id] = entity
		storage.LocationCount++
		//storage.LocationBytesDb[entity.Id] = entity.GetBytes()
	case "visits":
		entity := &models.Visit{}

		err := json.Unmarshal(ctx.PostBody, &params)
		if err != nil {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		if !entity.ValidateParams(params, "insert") {
			ctx.Error(fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
		storage.VisitDb[entity.Id] = entity
		storage.VisitCount++
		//storage.VisitBytesDb[entity.Id] = entity.GetBytes()

		user := storage.UserDb[entity.User]
		location := storage.LocationDb[entity.Location]

		entity.User_model = user
		entity.Location_model = location

		user.Visits = append(user.Visits, entity)
		location.Visits = append(location.Visits, entity)

		sort.Sort(models.VisitByDateAsc(user.Visits))
	}

	ctx.Write(EmptyJson)
}

