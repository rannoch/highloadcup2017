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

	CtxPool    sync.Pool
	ReaderPool sync.Pool
	WriterPool sync.Pool
	BytePool   sync.Pool
}

func New(port string) *TcpServer {
	fmt.Println("Creating server...")

	tcp_server := &TcpServer{
		Port: port,
	}

	return tcp_server
}

func (server *TcpServer) Listen() {
	fmt.Println("Launching server...")

	/*workerPool := WorkerPool{
		JobChan:make(chan net.Conn),
		WorkerFunc:server.ServeConn,
	}
	workerPool.Start()
*/
	listener, err := net.Listen("tcp4", ":"+server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()

		if err != nil {
			connection.Close()
			continue
		}

		go server.ServeConn(connection)
	}
}

func (server *TcpServer) ServeConn(c net.Conn) {
	hlcupCtx := acquireCtx(server, c)
	//hlcupCtx.ConnId = connId

	hlcupCtx.Handle()
}

func acquireReader(ctx *HlcupCtx) *bufio.Reader {
	v := ctx.Server.ReaderPool.Get()
	if v == nil {
		return bufio.NewReaderSize(ctx.Connection, 4096)
	}
	r := v.(*bufio.Reader)
	r.Reset(ctx.Connection)
	return r
}

func releaseReader(s *TcpServer, r *bufio.Reader) {
	s.ReaderPool.Put(r)
}

func acquireWriter(ctx *HlcupCtx) *bufio.Writer {
	v := ctx.Server.WriterPool.Get()
	if v == nil {
		return new(bufio.Writer)
	}
	w := v.(*bufio.Writer)
	w.Reset(ctx.Connection)
	return w
}

func releaseWriter(s *TcpServer, w *bufio.Writer) {
	s.WriterPool.Put(w)
}

func acquireCtx(s *TcpServer, c net.Conn) *HlcupCtx {
	v := s.CtxPool.Get()
	if v == nil {
		return NewCtx(s, c)
	}
	ctx := v.(*HlcupCtx)
	ctx.ResponseFullBuffer = acquireResponseFullBodyBuffer(s)
	ctx.ResponseBodyBuffer = acquireResponseBodyBuffer(s)

	ctx.Connection = c
	return ctx
}

func releaseCtx(s *TcpServer, ctx *HlcupCtx) {
	s.CtxPool.Put(ctx)
}

func acquireResponseBodyBuffer(s *TcpServer) *bytes.Buffer {
	v := s.ResponseBodyBufferPool.Get()
	if v == nil {
		return new(bytes.Buffer)
	}

	return v.(*bytes.Buffer)
}

func releaseResponseBodyBuffer(s *TcpServer, b *bytes.Buffer) {
	s.ResponseBodyBufferPool.Put(b)
}

func acquireResponseFullBodyBuffer(s *TcpServer) *bytes.Buffer {
	v := s.ResponseFullBufferPool.Get()
	if v == nil {
		return new(bytes.Buffer)
	}

	return v.(*bytes.Buffer)
}

func releaseResponseFullBufferPool(s *TcpServer, b *bytes.Buffer) {
	s.ResponseFullBufferPool.Put(b)
}

func acquireBytes(s *TcpServer) []byte {
	v := s.BytePool.Get()

	if v == nil {
		return make([]byte, 1024)
	}

	return v.([]byte)
}

func releaseBytes(s *TcpServer, b []byte) {
	s.BytePool.Put(b)
}