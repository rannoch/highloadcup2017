package server

import (
	"fmt"
	"net"
	"github.com/valyala/fasthttp"
	"bufio"
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

type HlcupCtx struct {
	ConnId           int64
	Method           []byte
	Url              []byte
	IsGet            bool
	IsPost           bool
	Connection       net.Conn
	ConnectionClosed bool
	KeepAlive        bool
	ResponseStatus   int

	Id           int64
	HasUrlParams bool
	UrlParams    []byte

	QueryArgs fasthttp.Args

	PostBody    []byte
	HasPostBody bool

	RequestBodyBuffer *bytes.Buffer

	ResponseBodyBuffer *bytes.Buffer
	ResponseFullBuffer *bytes.Buffer
}

func NewCtx(connection net.Conn) *HlcupCtx {
	responseBuffer := bufPool.Get().(*bytes.Buffer)
	responseBuffer.Reset()

	responseFullBuffer := bufPool.Get().(*bytes.Buffer)
	responseFullBuffer.Reset()

	return &HlcupCtx{
		Connection:         connection,
		ResponseBodyBuffer: responseBuffer,
		ResponseFullBuffer: responseFullBuffer,
		ResponseStatus:     200,
	}
}

func (hlcupRequest *HlcupCtx) ResetParams() {
	hlcupRequest.ResponseStatus = 200

	hlcupRequest.ResponseBodyBuffer.Reset()
	hlcupRequest.ResponseFullBuffer.Reset()
	hlcupRequest.QueryArgs.Reset()
}

func (hlcupRequest *HlcupCtx) Handle(handlerFunc func(ctx *HlcupCtx) (err error)) {
	buf := make([]byte, 1024)
	//net.TCPConn.SetKeepAlive(true)

	reader := bufio.NewReader(hlcupRequest.Connection)

	for {
		reader.Reset(hlcupRequest.Connection)
		n, err := reader.Read(buf)

		//fmt.Println(hlcupRequest.ConnId)

		if err != nil {
			hlcupRequest.Close()
			return
		}

		hlcupRequest.ResetParams()
		err = hlcupRequest.Parse(buf, n)

		handlerFunc(hlcupRequest)

		if hlcupRequest.ConnectionClosed {
			return
		}
	}
}

func (hlcupRequest *HlcupCtx) Parse(body []byte, n int) (err error) {
	if err != nil {
		//fmt.Println("Error reading:", err.Error())
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
		postBodyIndex := bytes.LastIndex(body, strCRLF)

		if postBodyIndex <= 0 {
			return
		}
		hlcupRequest.PostBody = body[postBodyIndex+2:n]
	}

	return
}

func (hlcupRequest *HlcupCtx) Error(status int) {
	hlcupRequest.ResponseStatus = status
}

func (hlcupRequest *HlcupCtx) Write(b []byte) {
	hlcupRequest.ResponseBodyBuffer.Write(b)
}

func (hlcupRequest *HlcupCtx) WriteString(s string) {
	hlcupRequest.ResponseBodyBuffer.WriteString(s)
}

func (hlcupRequest *HlcupCtx) Close() {
	hlcupRequest.ConnectionClosed = true

	hlcupRequest.Connection.Close()
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
			hlcupRequest.ResponseFullBuffer.Write(strContentLength)
			hlcupRequest.ResponseFullBuffer.Write(strColonSpace)
			hlcupRequest.ResponseFullBuffer.Write(strZero)
			hlcupRequest.ResponseFullBuffer.Write(strCRLF)
			hlcupRequest.ResponseFullBuffer.Write(strCRLF)

			hlcupRequest.Connection.Write(hlcupRequest.ResponseFullBuffer.Bytes())
			bufPool.Put(hlcupRequest.ResponseFullBuffer)
		} else {
			switch hlcupRequest.ResponseStatus {
			case 404:
				hlcupRequest.ResponseFullBuffer.Write(str404Closed)
			case 400:
				hlcupRequest.ResponseFullBuffer.Write(str400Closed)
			default:
				hlcupRequest.ResponseFullBuffer.Write(str404Closed)
			}

			hlcupRequest.ResponseFullBuffer.Write(strContentLength)
			hlcupRequest.ResponseFullBuffer.Write(strColonSpace)
			hlcupRequest.ResponseFullBuffer.Write(strZero)
			hlcupRequest.ResponseFullBuffer.Write(strCRLF)
			hlcupRequest.ResponseFullBuffer.Write(strCRLF)

			hlcupRequest.Connection.Write(hlcupRequest.ResponseFullBuffer.Bytes())
			bufPool.Put(hlcupRequest.ResponseFullBuffer)

			hlcupRequest.Close()
		}

		return
	}

	if hlcupRequest.KeepAlive {
		hlcupRequest.ResponseFullBuffer.Write(str200KeepAlive)
	} else {
		hlcupRequest.ResponseFullBuffer.Write(str200Closed)
	}

	hlcupRequest.ResponseFullBuffer.Write(strContentLength)
	hlcupRequest.ResponseFullBuffer.Write(strColonSpace)
	hlcupRequest.ResponseFullBuffer.WriteString(fmt.Sprintf("%d", len(hlcupRequest.ResponseBodyBuffer.Bytes())))
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)

	hlcupRequest.ResponseFullBuffer.Write(hlcupRequest.ResponseBodyBuffer.Bytes())
	bufPool.Put(hlcupRequest.ResponseBodyBuffer)

	hlcupRequest.Connection.Write(hlcupRequest.ResponseFullBuffer.Bytes())

	bufPool.Put(hlcupRequest.ResponseFullBuffer)

	if !hlcupRequest.KeepAlive {
		hlcupRequest.Close()
	}
}
