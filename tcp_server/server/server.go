package server

import (
	"fmt"
	"net"
	"os"
)

type TcpServer struct {
	Port string
	HandleFunc func(hlcupCtx *HlcupCtx) (err error)
}

func New(port string) *TcpServer {
	fmt.Println("Creating server...")

	tcp_server := &TcpServer{
		Port : port,
	}

	return tcp_server
}

func (server *TcpServer) Listen(){
	fmt.Println("Launching server...")

	listener, err := net.Listen("tcp4", ":" + server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	var connId int64

	for {
		connection, err := listener.Accept()

		connId++

		if err != nil {
			fmt.Println(err.Error())
			connection.Close()
			continue
		}

		hlcupCtx := NewCtx(connection)
		hlcupCtx.ConnId = connId

		go hlcupCtx.Handle(server.HandleFunc)
	}
}