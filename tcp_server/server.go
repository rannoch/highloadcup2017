package main

import (
	//"bufio"
	//"net/textproto"
	"log"
	"os"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"github.com/valyala/fasthttp"
	"github.com/rannoch/highloadcup2017/tcp_server/handlers"
	"github.com/rannoch/highloadcup2017/tcp_server/server"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("not enough args")
	}

	//storage.Init()

	storage.InitMemoryMap()
	LoadData(os.Args[2])

	fmt.Println("Launching server...")

	// listen on all interfaces
	listener, _ := net.Listen("tcp4", ":8086")
	defer listener.Close()

	// run loop forever (or until ctrl-c)
	for {
		// accept connection on port
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println(err.Error())
			connection.Close()
			continue
			//panic("listener accept error")
		}

		hlcupCtx := server.HlcupCtx{
			Connection:connection,
			HasUrlParams:true,
		}
		hlcupCtx.TryParse()

		go Handler(&hlcupCtx)
	}
}

func Handler(hlcupCtx *server.HlcupCtx) {
	path := hlcupCtx.Url

	if hlcupCtx.IsPost {
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
			return
		}

		params := bytes.Split(path, []byte("/"))
		if len(params) != 2 {
			hlcupCtx.Error(400)
			return
		}

		entity := params[0]
		idValue := params[1]

		if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
			hlcupCtx.Error(fasthttp.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(string(idValue[:]))

		if err != nil || id < 0{
			hlcupCtx.Error(fasthttp.StatusNotFound)
			return
		}

		//POST /<entity>/<id> на обновление
		handlers.EntityUpdateHandler(hlcupCtx, int64(id), entity)
		return
	}

	if hlcupCtx.IsGet {
		params := bytes.Split(path, []byte("/"))
		if len(params) < 2 || len(params) > 3 {
			hlcupCtx.Error(fasthttp.StatusBadRequest)
			return
		}

		entity := params[0]
		idValue := params[1]

		//fmt.Println(string(entity))
		//fmt.Println(string(idValue))

		if !(bytes.Equal(entity, handlers.UsersBytes) || bytes.Equal(entity, handlers.LocationsBytes) || bytes.Equal(entity, handlers.VisitsBytes)) {
			hlcupCtx.Error(fasthttp.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(string(idValue[:]))
		if err != nil || id < 0{
			hlcupCtx.Error(fasthttp.StatusNotFound)
			return
		}

		if len(params) == 3 {
			if bytes.Equal(entity, handlers.UsersBytes) && bytes.Equal(params[2], handlers.VisitsBytes) {
				handlers.UsersVisitsHandler(hlcupCtx, int64(id))
				return
			}
			if bytes.Equal(entity, handlers.LocationsBytes) && bytes.Equal(params[2], handlers.AvgBytes) {
				handlers.LocationsAvgHandler(hlcupCtx, int64(id))
				return
			}

			hlcupCtx.Error(fasthttp.StatusBadRequest)
			return
		}

		handlers.EntitySelectHandler(hlcupCtx, int64(id), entity)
		return
	}

	hlcupCtx.Error(fasthttp.StatusNotFound)
}



