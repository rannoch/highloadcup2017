package main

import (
	//"bufio"
	//"net/textproto"
	"log"
	"os"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
	"github.com/valyala/fasthttp"
	"bytes"
	"strconv"
	"github.com/rannoch/highloadcup2017/tcp_server/handlers"
)

var strSlash = []byte("/")

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	//storage.Init()

	storage.InitMemoryMap()
	LoadData(os.Args[2])

	tcpServer := server.New(os.Args[1])
	tcpServer.HandleFunc = HandleFunc

	tcpServer.Listen()
}

func HandleFunc(hlcupCtx *server.HlcupCtx) (err error) {
	path := hlcupCtx.Url

	if hlcupCtx.IsPost {
		if !hlcupCtx.HasPostBody {
			hlcupCtx.Error(400)
			hlcupCtx.SendResponse()
			return
		}

		//POST /<entity>/new на создание
		if bytes.Equal(path, handlers.UsersNewBytes) || bytes.Equal(path, handlers.LocationsNewBytes) || bytes.Equal(path, handlers.VisitsNewBytes) {
			switch {
			case bytes.Equal(path, handlers.UsersNewBytes):
				handlers.EntitityNewHandler(hlcupCtx, "users")
			case bytes.Equal(path, handlers.LocationsNewBytes):
				handlers.EntitityNewHandler(hlcupCtx, "locations")
			case bytes.Equal(path, handlers.VisitsNewBytes):
				handlers.EntitityNewHandler(hlcupCtx, "visits")
			}

			hlcupCtx.SendResponse()
			return
		}

		params := bytes.Split(path, strSlash)
		if len(params) != 2 {
			hlcupCtx.Error(400)
			hlcupCtx.SendResponse()
			return
		}

		entity := params[0]
		idValue := params[1]

		if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
			hlcupCtx.Error(fasthttp.StatusBadRequest)
			hlcupCtx.SendResponse()
			return
		}

		//id, err := strconv.Atoi(string(idValue[:]))
		id, err := strconv.ParseInt(string(idValue[:]), 10, 0)

		if err != nil || id < 0 {
			hlcupCtx.Error(fasthttp.StatusNotFound)
			hlcupCtx.SendResponse()
			return err
		}

		//POST /<entity>/<id> на обновление
		handlers.EntityUpdateHandler(hlcupCtx, id, entity)
		hlcupCtx.SendResponse()
		return err
	}

	if hlcupCtx.IsGet {
		params := bytes.Split(path, strSlash)
		if len(params) < 2 || len(params) > 3 {
			hlcupCtx.Error(fasthttp.StatusBadRequest)
			hlcupCtx.SendResponse()
			return
		}

		entity := params[0]
		idValue := params[1]

		if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
			hlcupCtx.Error(fasthttp.StatusNotFound)
			hlcupCtx.SendResponse()
			return
		}

		id, err := strconv.ParseInt(string(idValue[:]), 10, 0)
		if err != nil || id < 0 {
			hlcupCtx.Error(fasthttp.StatusNotFound)
			hlcupCtx.SendResponse()
			return err
		}

		if len(params) == 3 {
			if bytes.Equal(entity, handlers.UsersBytes) && bytes.Equal(params[2], handlers.VisitsBytes) {
				handlers.UsersVisitsHandler(hlcupCtx, id)
				hlcupCtx.SendResponse()
				return err
			}
			if bytes.Equal(entity, handlers.LocationsBytes) && bytes.Equal(params[2], handlers.AvgBytes) {
				handlers.LocationsAvgHandler(hlcupCtx, id)
				hlcupCtx.SendResponse()
				return err
			}

			hlcupCtx.Error(fasthttp.StatusBadRequest)
			hlcupCtx.SendResponse()
			return err
		}

		handlers.EntitySelectHandler(hlcupCtx, id, entity)
		hlcupCtx.SendResponse()
		return err
	}

	hlcupCtx.Error(fasthttp.StatusNotFound)
	hlcupCtx.SendResponse()

	return
}
