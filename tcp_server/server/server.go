package server

import (
	"fmt"
	"net"
	"os"
	"sync"
	"bytes"
	"bufio"
)

type TcpServer struct {
	Port       string
	HandleFunc func(hlcupCtx *HlcupCtx) (err error)

	ResponseBodyBufferPool sync.Pool
	ResponseFullBufferPool sync.Pool

	ctxPool    sync.Pool
	ReaderPool sync.Pool
	WriterPool sync.Pool
	BytePool sync.Pool
}

func New(port string) *TcpServer {
	fmt.Println("Creating server...")

	tcp_server := &TcpServer{
		Port: port,
	}

	tcp_server.ResponseBodyBufferPool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(bytes.Buffer)
		},
	}

	tcp_server.ResponseFullBufferPool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(bytes.Buffer)
		},
	}

	tcp_server.ReaderPool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(bufio.Reader)
		},
	}

	tcp_server.WriterPool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return new(bufio.Writer)
		},
	}

	tcp_server.BytePool = sync.Pool{
		New: func() interface{} {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return make([]byte, 1024)
		},
	}

	return tcp_server
}

func (server *TcpServer) Listen() {
	fmt.Println("Launching server...")

	listener, err := net.Listen("tcp4", ":"+server.Port)
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

		go func() {
			hlcupCtx := NewCtx(server, connection)
			hlcupCtx.ConnId = connId

			hlcupCtx.Handle()
		}()
	}
}
