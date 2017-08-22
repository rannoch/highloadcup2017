package handlers

import (
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"github.com/rannoch/highloadcup2017/memory/models"
	"encoding/json"
	"sort"
	"bytes"
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
		if id > storage.UserCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.UserDb[id]

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case "locations":
		if id > storage.LocationCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.LocationDb[id]

		if !entity.ValidateParams(params, "update") {
			ctx.Error("", fasthttp.StatusBadRequest)
			return
		}

		entity.SetParams(params)
	case "visits":
		if id > storage.VisitCount {
			ctx.Error("", fasthttp.StatusNotFound)
			return
		}

		entity := storage.VisitDb[id]

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
			sort.Sort(models.VisitByDateAsc(userUpdated.Visits))
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

		_, ok = params["visited_at"]
		if ok {
			sort.Sort(models.VisitByDateAsc(entity.User_model.Visits))
		}
	}

	buffer := bufPool.Get().(*bytes.Buffer)
	buffer.Reset()
	buffer.WriteString("{}")

	ctx.Write(buffer.Bytes())
	bufPool.Put(buffer)
}