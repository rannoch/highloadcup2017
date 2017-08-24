package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"github.com/rannoch/highloadcup2017/memory/models"
	"encoding/json"
	"sort"
	"bytes"
)

func EntityUpdateHandler(ctx *fasthttp.RequestCtx, id int64, entityValue []byte) {
	ctx.SetContentType("application/json;charset=utf-8")

	defer ctx.SetConnectionClose()

	switch {
	case bytes.Equal(entityValue, UsersBytes):
		if id > storage.UserCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		// check params
		var params map[string]interface{}
		err := json.Unmarshal(ctx.PostBody(), &params)

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity := storage.UserDb[id]

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)

		storage.UserBytesDb[entity.Id] = entity.GetBytes()
	case bytes.Equal(entityValue, LocationsBytes):
		if id > storage.LocationCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		// check params
		var params map[string]interface{}
		err := json.Unmarshal(ctx.PostBody(), &params)

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity := storage.LocationDb[id]

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)

		storage.LocationBytesDb[entity.Id] = entity.GetBytes()
	case bytes.Equal(entityValue, VisitsBytes):
		if id > storage.VisitCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		// check params
		var params map[string]interface{}
		err := json.Unmarshal(ctx.PostBody(), &params)

		if err != nil {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity := storage.VisitDb[id]

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		userParam, ok := params["user"]
		if ok && userParam != entity.User {
			var userIdOld int64 = entity.User
			var userIdUpdated int64
			switch userParam.(type) {
			case int64:
				userIdUpdated = userParam.(int64)
			case int32:
				userIdUpdated = int64(userParam.(int32))
			case float32:
				userIdUpdated = int64(userParam.(float32))
			case float64:
				userIdUpdated = int64(userParam.(float64))
			}

			userUpdated := storage.UserDb[userIdUpdated]
			userOld := storage.UserDb[userIdOld]

			entity.User_model = userUpdated

			// удаляю визит из старого пользователя
			for i, visit := range userOld.Visits {
				if visit.Id == entity.Id {
					userOld.Visits = append(userOld.Visits[:i], userOld.Visits[i+1:]...)
					break
				}
			}
			// добавляю в нового
			userUpdated.Visits = append(userUpdated.Visits, entity)
			sort.Sort(models.VisitByDateAsc(userUpdated.Visits))
		}

		locationParam, ok := params["location"]
		if ok && locationParam != entity.Location {
			var locationIdOld int64 = entity.Location
			var locationIdUpdated int64
			switch locationParam.(type) {
			case int64:
				locationIdUpdated = locationParam.(int64)
			case int32:
				locationIdUpdated = int64(locationParam.(int32))
			case float32:
				locationIdUpdated = int64(locationParam.(float32))
			case float64:
				locationIdUpdated = int64(locationParam.(float64))
			}

			locationUpdated := storage.LocationDb[locationIdUpdated]
			locationOld := storage.LocationDb[locationIdOld]

			entity.Location_model = locationUpdated

			// удаляю визит из старого пользователя
			for i, visit := range locationOld.Visits {
				if visit.Id == entity.Id {
					locationOld.Visits = append(locationOld.Visits[:i], locationOld.Visits[i+1:]...)
					break
				}
			}
			// добавляю в нового
			locationUpdated.Visits = append(locationUpdated.Visits, entity)
		}

		entity.SetParams(params)
		storage.VisitBytesDb[entity.Id] = entity.GetBytes()

		_, ok = params["visited_at"]
		if ok {
			sort.Sort(models.VisitByDateAsc(entity.User_model.Visits))
		}
	}

	ctx.SetBody(EmptyJson)
}