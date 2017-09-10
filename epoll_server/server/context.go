package server

import (
	"github.com/valyala/fasthttp"
	"bytes"
	"strconv"
	"syscall"
)

type HlcupCtx struct {
	ConnId           int64
	Method           []byte
	Url              []byte
	IsGet            bool
	IsPost           bool
	Fd               int32
	ConnectionClosed bool
	KeepAlive        bool
	ResponseStatus   int

	Id           int64
	HasUrlParams bool
	UrlParams    []byte

	QueryArgs fasthttp.Args

	PostBody    []byte
	HasPostBody bool

	ResponseBodyBuffer *bytes.Buffer
	ResponseFullBuffer *bytes.Buffer

	Server *Epoll_server

	//Header fasthttp.RequestHeader
}

func NewCtx(server *Epoll_server, fd int32) *HlcupCtx {
	return &HlcupCtx{
		Fd:                 fd,
		ResponseStatus:     200,
		Server:             server,
		ResponseBodyBuffer: acquireResponseBodyBuffer(server),
		ResponseFullBuffer: acquireResponseFullBodyBuffer(server),
	}
}

func (hlcupRequest *HlcupCtx) ResetParams() {
	hlcupRequest.Method = hlcupRequest.Method[:0]
	hlcupRequest.ResponseStatus = 200
	hlcupRequest.KeepAlive = false

	hlcupRequest.ResponseBodyBuffer.Reset()
	hlcupRequest.ResponseFullBuffer.Reset()

	hlcupRequest.QueryArgs.Reset()

	hlcupRequest.IsPost = false
	hlcupRequest.IsGet = false

	hlcupRequest.HasUrlParams = false
	hlcupRequest.UrlParams = hlcupRequest.UrlParams[:0]
	hlcupRequest.Url = hlcupRequest.Url[:0]

	hlcupRequest.PostBody = hlcupRequest.PostBody[:0]
	hlcupRequest.HasPostBody = false
}

func (hlcupRequest *HlcupCtx) Handle() {
	hlcupRequest.ResetParams()

	buf := acquireBytes(hlcupRequest.Server)

	n, err := syscall.Read(int(hlcupRequest.Fd), buf[:])

	if err != nil {
		releaseBytes(hlcupRequest.Server, buf)

		releaseResponseBodyBuffer(hlcupRequest.Server, hlcupRequest.ResponseBodyBuffer)
		releaseResponseFullBufferPool(hlcupRequest.Server, hlcupRequest.ResponseFullBuffer)
		releaseCtx(hlcupRequest.Server, hlcupRequest)
		return
	}

	err = hlcupRequest.Parse(buf, n)

	releaseBytes(hlcupRequest.Server, buf)

	hlcupRequest.Server.HandleFunc(hlcupRequest)

	if !hlcupRequest.KeepAlive {
		hlcupRequest.Close()
	}

	releaseResponseBodyBuffer(hlcupRequest.Server, hlcupRequest.ResponseBodyBuffer)
	releaseResponseFullBufferPool(hlcupRequest.Server, hlcupRequest.ResponseFullBuffer)
	releaseCtx(hlcupRequest.Server, hlcupRequest)
}

func (hlcupRequest *HlcupCtx) Parse(body []byte, n int) (err error) {
	if err != nil {
		return
	}

	// method
	methodIndex := bytes.IndexByte(body, ' ')
	if methodIndex <= 0 {
		return
	}
	hlcupRequest.Method = body[:methodIndex]

	if bytes.Equal(hlcupRequest.Method, strGet) {
		hlcupRequest.IsGet = true
		hlcupRequest.KeepAlive = true
	}

	if bytes.Equal(hlcupRequest.Method, strPost) {
		hlcupRequest.IsPost = true
	}

	// url
	urlIndex := bytes.IndexByte(body[methodIndex+1:], '?')

	if urlIndex <= 0 {
		hlcupRequest.HasUrlParams = false
		hlcupRequest.UrlParams = hlcupRequest.UrlParams[:0]

		urlIndex = bytes.IndexByte(body[methodIndex+1:], ' ')

		if urlIndex <= 0 {
			return
		}
	} else {
		hlcupRequest.HasUrlParams = true
	}

	hlcupRequest.Url = body[methodIndex+2:methodIndex+1+urlIndex]

	// params
	if hlcupRequest.HasUrlParams {
		paramsIndex := bytes.IndexByte(body[methodIndex+urlIndex+2:], ' ')
		if paramsIndex <= 0 {
			return
		}

		hlcupRequest.UrlParams = body[methodIndex+urlIndex+2: methodIndex+urlIndex+2+paramsIndex]
	}

	// postBody
	if hlcupRequest.IsPost {
		postBodyIndex := bytes.Index(body, append(strCRLF[:], strCRLF[:]...))

		if postBodyIndex <= 0 || postBodyIndex >= n {
			hlcupRequest.HasPostBody = false
			return
		}

		hlcupRequest.HasPostBody = true
		hlcupRequest.PostBody = body[postBodyIndex+2:n]
	}

	return
}

func (hlcupRequest *HlcupCtx) Error(status int) {
	hlcupRequest.ResponseStatus = status
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	//log.Println("error")
}

func (hlcupRequest *HlcupCtx) Write(b []byte) {
	hlcupRequest.ResponseBodyBuffer.Write(b)
}

func (hlcupRequest *HlcupCtx) WriteString(s string) {
	hlcupRequest.ResponseBodyBuffer.WriteString(s)
}

func (hlcupRequest *HlcupCtx) Close() {
	hlcupRequest.ConnectionClosed = true
	syscall.Close(int(hlcupRequest.Fd))
	//hlcupRequest.Fd.Close()
}

func (hlcupRequest *HlcupCtx) ParseParams() (params map[string]interface{}) {
	params = map[string]interface{}{}

	if !hlcupRequest.HasUrlParams {
		return
	}

	hlcupRequest.QueryArgs.ParseBytes(hlcupRequest.UrlParams)

	return
}

func (hlcupRequest *HlcupCtx) SendResponse() {
	//writer := acquireWriter(hlcupRequest)
	//writer.Reset(hlcupRequest.Fd)

	if hlcupRequest.ResponseStatus != 200 {
		if hlcupRequest.KeepAlive {
			switch hlcupRequest.ResponseStatus {
			case 404:
				hlcupRequest.ResponseFullBuffer.Write(str404KeepAlive)
			case 400:
				hlcupRequest.ResponseFullBuffer.Write(str400KeepAlive)
			default:
				hlcupRequest.ResponseFullBuffer.Write(str404KeepAlive)
			}
		} else {
			switch hlcupRequest.ResponseStatus {
			case 404:
				hlcupRequest.ResponseFullBuffer.Write(str404Closed)
			case 400:
				hlcupRequest.ResponseFullBuffer.Write(str400Closed)
			default:
				hlcupRequest.ResponseFullBuffer.Write(str404Closed)
			}
		}

		hlcupRequest.ResponseFullBuffer.Write(strContentLength)
		hlcupRequest.ResponseFullBuffer.Write(strColonSpace)
		hlcupRequest.ResponseFullBuffer.Write(strZero)
		hlcupRequest.ResponseFullBuffer.Write(strCRLF)
		hlcupRequest.ResponseFullBuffer.Write(strCRLF)

		//writer.Write(hlcupRequest.ResponseFullBuffer.Bytes())
		//writer.Flush()
		//releaseWriter(hlcupRequest.Server, writer)

		syscall.Write(int(hlcupRequest.Fd), hlcupRequest.ResponseFullBuffer.Bytes())

		return
	}

	if hlcupRequest.KeepAlive {
		hlcupRequest.ResponseFullBuffer.Write(str200KeepAlive)
	} else {
		hlcupRequest.ResponseFullBuffer.Write(str200Closed)
	}

	hlcupRequest.ResponseFullBuffer.Write(strContentLength)
	hlcupRequest.ResponseFullBuffer.Write(strColonSpace)

	hlcupRequest.ResponseFullBuffer.WriteString(strconv.Itoa(len(hlcupRequest.ResponseBodyBuffer.Bytes())))
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)

	hlcupRequest.ResponseFullBuffer.Write(hlcupRequest.ResponseBodyBuffer.Bytes())

	syscall.Write(int(hlcupRequest.Fd), hlcupRequest.ResponseFullBuffer.Bytes())

	//writer.Write(hlcupRequest.ResponseFullBuffer.Bytes())
	//writer.Flush()

	//fmt.Println(hlcupRequest.ResponseFullBuffer.String())
	//releaseWriter(hlcupRequest.Server, writer)
}
