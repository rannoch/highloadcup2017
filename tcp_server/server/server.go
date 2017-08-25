package server

import (
	"net"
	"bufio"
	"bytes"
	//"fmt"
	"github.com/valyala/fasthttp"
	//"fmt"
	"fmt"
	//"io"
	//"time"
	"time"
)

type HlcupCtx struct {
	Method []byte
	Url    []byte
	IsGet  bool
	IsPost bool
	Connection net.Conn
	ConnectionClosed bool

	Id     int64
	HasUrlParams bool
	UrlParams []byte

	QueryArgs fasthttp.Args

	PostBody []byte
	HasPostBody bool
}

func (hlcupRequest *HlcupCtx) Handle(handlerFunc func(ctx *HlcupCtx) (err error)) {
	reader := bufio.NewReader(hlcupRequest.Connection)

	buf := make([]byte, 2048)

	//net.TCPConn.SetKeepAlive(true)

	for {
		//hlcupRequest.Connection.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

		n, err := reader.Read(buf)

		if err != nil {
			hlcupRequest.Close()
			return
		}

		err = hlcupRequest.Parse(buf, n)

		handlerFunc(hlcupRequest)

		if hlcupRequest.ConnectionClosed {
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func (hlcupRequest *HlcupCtx) Parse(body []byte, n int) (err error) {
	//fmt.Println("START ---------")
	//fmt.Println(string(body))
	//fmt.Println("END ---------")
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		//err = io.EOF
		return
	}

	// method
	methodIndex := bytes.IndexByte(body, ' ')
	if methodIndex <= 0 {
		return
	}
	hlcupRequest.Method = body[:methodIndex]

	//fmt.Println(string(hlcupRequest.Method))

	if bytes.Equal(hlcupRequest.Method, strGet) {
		hlcupRequest.IsGet = true
	}

	if bytes.Equal(hlcupRequest.Method, strPost) {
		hlcupRequest.IsPost = true
	}

	// url
	urlIndex := bytes.IndexByte(body[methodIndex + 1:], '?')
	if urlIndex <= 0 {
		hlcupRequest.HasUrlParams = false

		urlIndex = bytes.IndexByte(body[methodIndex + 1:], ' ')

		if urlIndex <= 0 {
			return
		}
	}

	hlcupRequest.Url = body[methodIndex + 2:methodIndex + 1 + urlIndex]

	//fmt.Println(string(hlcupRequest.Url))

	// params
	if hlcupRequest.HasUrlParams {
		paramsIndex := bytes.IndexByte(body[methodIndex + urlIndex + 2:], ' ')
		if paramsIndex <= 0 {
			return
		}

		hlcupRequest.UrlParams = body[methodIndex + urlIndex + 2: methodIndex + urlIndex + 2 + paramsIndex]
	}

	// postBody
	if hlcupRequest.IsPost {
		postBodyIndex := bytes.LastIndex(body, strCRLF)

		if postBodyIndex <= 0 {
			return
		}
		hlcupRequest.PostBody = body[postBodyIndex + 2 :n]

		//fmt.Println("." + string(hlcupRequest.PostBody) + ".")
	}

	return
}

func (hlcupRequest *HlcupCtx) Error(status int) {
	if hlcupRequest.IsGet {
		switch status {
		case 404:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Keep-Alive\n\n"))
		case 400:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 400 OK\ncontent-type:application/json;charset=utf-8;Connection: Keep-Alive\n\n"))
		default:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Keep-Alive\n\n"))
		}
	} else {
		switch status {
		case 404:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
		case 400:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 400 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
		default:
			hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
		}
	}

	hlcupRequest.Close()
}

func (hlcupRequest *HlcupCtx) SetBody(b []byte) {
	if hlcupRequest.IsGet {
		hlcupRequest.Connection.Write([]byte("HTTP/1.1 200 OK\ncontent-type:application/json;charset=utf-8;\nConnection: Keep-Alive;\nTransfer-Encoding: chunked"))
	} else {
		hlcupRequest.Connection.Write([]byte("HTTP/1.1 200 OK\ncontent-type:application/json;charset=utf-8;\nConnection: Closed;Transfer-Encoding: chunked"))
	}

	hlcupRequest.Connection.Write(strCRLF)
	hlcupRequest.Connection.Write(strCRLF)

	length := []byte(fmt.Sprintf("%x", len(string(b[:]))))
	//
	//fmt.Println(string(b[:]))
	//fmt.Println(len(b[:]))
	//fmt.Println(len(string(b[:])))
	hlcupRequest.Connection.Write(length)
	hlcupRequest.Connection.Write(strCRLF)
	hlcupRequest.Connection.Write(b)
	hlcupRequest.Connection.Write(strCRLF)
	hlcupRequest.Connection.Write(strZero)
	hlcupRequest.Connection.Write(strCRLF)
	hlcupRequest.Connection.Write(strCRLF)
}

func (hlcupRequest *HlcupCtx) Write(b []byte) {
	hlcupRequest.Connection.Write(b)
}

func (hlcupRequest *HlcupCtx) WriteString(s string) {
	hlcupRequest.Connection.Write([]byte(s))
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