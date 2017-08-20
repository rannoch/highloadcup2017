package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"strconv"
	"strings"
	"github.com/rannoch/highloadcup2017/memory/models"
	"encoding/json"
)

func EntityUpdateHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json;charset=utf-8")

	var id int
	var entityValue string
	var params map[string]interface{}

	if ctx.UserValue("id").(string) == "new" {
		EntitityNewHandler(ctx)
		return
	}

	id, err := strconv.Atoi(ctx.UserValue("id").(string))

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	// check params
	err = json.Unmarshal(ctx.PostBody(), &params)

	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch {
	case strings.Contains(entityValue, "user"):
		e, ok := storage.Db["user"][int32(id)]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := e.(*models.User)

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case strings.Contains(entityValue, "location"):
		e, ok := storage.Db["location"][int32(id)]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := e.(*models.Location)

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case strings.Contains(entityValue, "visit"):
		e, ok := storage.Db["visit"][int32(id)]
		if !ok {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := e.(*models.Visit)

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

			userUpdated := storage.Db["user"][userIdUpdated].(*models.User)
			userOld := storage.Db["user"][userIdOld].(*models.User)

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

			locationUpdated := storage.Db["location"][locationIdUpdated].(*models.Location)
			locationOld := storage.Db["location"][locationIdOld].(*models.Location)

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

	entityValue, ok := ctx.UserValue("entity").(string)

	if !ok {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}

	switch {
	case strings.Contains(entityValue, "user"):
		entity := &models.User{}
		// check params
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

		storage.Db["user"][entity.Id] = entity
	case strings.Contains(entityValue, "location"):
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

		storage.Db["location"][entity.Id] = entity
	case strings.Contains(entityValue, "visit"):
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
		storage.Db["visit"][entity.Id] = entity

		user := storage.Db["user"][entity.User].(*models.User)
		location := storage.Db["location"][entity.Location].(*models.Location)

		entity.User_model = user
		entity.Location_model = location

		user.Visits = append(user.Visits, entity)
		location.Visits = append(location.Visits, entity)
	}

	ctx.SetBody([]byte("{}"))
}
