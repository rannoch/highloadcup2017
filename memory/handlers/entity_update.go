package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"github.com/rannoch/highloadcup2017/memory/models"
	"encoding/json"
	"github.com/antonholmquist/jason"
)

func EntityUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int32
	var entityValue string
	var params map[string]interface{}

	defer ctx.SetConnectionClose()

	id, _ = ctx.UserValue("id").(int32)

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	// check params
	err := json.Unmarshal(ctx.PostBody(), &params)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch entityValue {
	case "users":
		entity, ok := storage.UserDb[id]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case "locations":
		entity, ok := storage.LocationDb[id]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case "visits":
		entity, ok := storage.VisitDb[id]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		userParam, ok := params["user"]
		if ok && userParam != entity.User {
			var userIdOld int32 = entity.User
			var userIdUpdated int32
			switch userParam.(type) {
			case int32:
				userIdUpdated = userParam.(int32)
			case float32:
				userIdUpdated = int32(userParam.(float32))
			case float64:
				userIdUpdated = int32(userParam.(float64))
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
		}

		locationParam, ok := params["location"]
		if ok && locationParam != entity.Location {
			var locationIdOld int32 = entity.Location
			var locationIdUpdated int32
			switch locationParam.(type) {
			case int32:
				locationIdUpdated = locationParam.(int32)
			case float32:
				locationIdUpdated = int32(locationParam.(float32))
			case float64:
				locationIdUpdated = int32(locationParam.(float64))
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
	}

	ctx.SetBody([]byte("{}"))
}

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

		user := storage.UserDb[entity.User]
		location := storage.LocationDb[entity.Location]

		entity.User_model = user
		entity.Location_model = location

		user.Visits = append(user.Visits, entity)
		location.Visits = append(location.Visits, entity)
	}

	ctx.SetBody([]byte("{}"))
}
