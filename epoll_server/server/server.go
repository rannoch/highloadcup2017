package server

import (
	"fmt"
	"net"
	"os"
	"sync"
	"bytes"
	"syscall"
)

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 32
)

type Epoll_server struct {
	Port       int
	HandleFunc func(hlcupCtx *HlcupCtx) (err error)

	ResponseBodyBufferPool sync.Pool
	ResponseFullBufferPool sync.Pool

	CtxPool    sync.Pool
	ReaderPool sync.Pool
	WriterPool sync.Pool
	BytePool   sync.Pool
}

func New(port int) *Epoll_server {
	fmt.Println("Creating server...")

	epoll_server := &Epoll_server{
		Port: port,
	}

	return epoll_server
}

func (server *Epoll_server) Listen() {
	fmt.Println("Launching server...")

	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent

	// socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer syscall.Close(socket)

	if err = syscall.SetNonblock(socket, true); err != nil {
		fmt.Println("setnonblock1: ", err)
		os.Exit(1)
	}

	addr := syscall.SockaddrInet4{Port: server.Port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())

	err = syscall.Bind(socket, &addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = syscall.Listen(socket, 1000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// epoll
	epoll, e := syscall.EpollCreate1(0)
	if e != nil {
		fmt.Println("epoll_create1: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epoll)

	event.Events = syscall.EPOLLIN
	event.Fd = int32(socket)

	if err = syscall.EpollCtl(epoll, syscall.EPOLL_CTL_ADD, socket, &event); err != nil {
		fmt.Println("epoll_ctl: ", err)
		os.Exit(1)
	}

	for {
		num_events, e := syscall.EpollWait(epoll, events[:], -1)
		if e != nil {
			fmt.Println("epoll_wait: ", e)
			break
		}

		for ev := 0; ev < num_events; ev++ {
			if int(events[ev].Fd) == socket {
				connFd, _, err := syscall.Accept(socket)
				if err != nil {
					fmt.Println("accept: ", err)
					continue
				}
				syscall.SetNonblock(socket, true)
				event.Events = syscall.EPOLLIN | EPOLLET
				event.Fd = int32(connFd)
				if err := syscall.EpollCtl(epoll, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
					fmt.Print("epoll_ctl: ", connFd, err)
					os.Exit(1)
				}
			} else {
				//go echo(int(server.Events[ev].Fd))

				go server.ServeConn(events[ev].Fd)
			}
		}
	}
}

func echo(socket int) {
	defer syscall.Close(socket)
	var buf [32 * 1024]byte
	for {
		nbytes, e := syscall.Read(socket, buf[:])
		if nbytes > 0 {
			fmt.Printf(">>> %s", buf)
			syscall.Write(socket, buf[:nbytes])
			fmt.Printf("<<< %s", buf)
		}
		if e != nil {
			break
		}
	}
}

func (server *Epoll_server) ServeConn(fd int32) {
	hlcupCtx := acquireCtx(server, fd)
	//hlcupCtx.ConnId = connId

	hlcupCtx.Handle()
}

/*func acquireReader(ctx *HlcupCtx) *bufio.Reader {
	v := ctx.Server.ReaderPool.Get()
	if v == nil {
		return bufio.NewReaderSize(ctx.Fd, 4096)
	}
	r := v.(*bufio.Reader)
	r.Reset(ctx.Fd)
	return r
}

func releaseReader(s *Epoll_server, r *bufio.Reader) {
	s.ReaderPool.Put(r)
}

func acquireWriter(ctx *HlcupCtx) *bufio.Writer {
	v := ctx.Server.WriterPool.Get()
	if v == nil {
		return new(bufio.Writer)
	}
	w := v.(*bufio.Writer)
	w.Reset(ctx.Fd)
	return w
}*/

/*func releaseWriter(s *Epoll_server, w *bufio.Writer) {
	s.WriterPool.Put(w)
}*/

func acquireCtx(s *Epoll_server, fd int32) *HlcupCtx {
	v := s.CtxPool.Get()
	if v == nil {
		return NewCtx(s, fd)
	}
	ctx := v.(*HlcupCtx)
	ctx.ResponseFullBuffer = acquireResponseFullBodyBuffer(s)
	ctx.ResponseBodyBuffer = acquireResponseBodyBuffer(s)

	ctx.Fd = fd
	return ctx
}

func releaseCtx(s *Epoll_server, ctx *HlcupCtx) {
	s.CtxPool.Put(ctx)
}

func acquireResponseBodyBuffer(s *Epoll_server) *bytes.Buffer {
	v := s.ResponseBodyBufferPool.Get()
	if v == nil {
		return new(bytes.Buffer)
	}

	return v.(*bytes.Buffer)
}

func releaseResponseBodyBuffer(s *Epoll_server, b *bytes.Buffer) {
	s.ResponseBodyBufferPool.Put(b)
}

func acquireResponseFullBodyBuffer(s *Epoll_server) *bytes.Buffer {
	v := s.ResponseFullBufferPool.Get()
	if v == nil {
		return new(bytes.Buffer)
	}

	return v.(*bytes.Buffer)
}

func releaseResponseFullBufferPool(s *Epoll_server, b *bytes.Buffer) {
	s.ResponseFullBufferPool.Put(b)
}

func acquireBytes(s *Epoll_server) []byte {
	v := s.BytePool.Get()

	if v == nil {
		return make([]byte, 1024)
	}

	return v.([]byte)
}

func releaseBytes(s *Epoll_server, b []byte) {
	s.BytePool.Put(b)
}