package server

import (
	"net"
	"bufio"
	"bytes"
	//"fmt"
	"github.com/valyala/fasthttp"
	//"fmt"
	"fmt"
	"io"
)

type HlcupCtx struct {
	Method []byte
	Url    []byte
	IsGet  bool
	IsPost bool
	Connection net.Conn

	Id     int64
	HasUrlParams bool
	UrlParams []byte

	QueryArgs fasthttp.Args

	PostBody []byte
	HasPostBody bool
}

func (hlcupRequest *HlcupCtx) TryParse() {
	for {
		err := hlcupRequest.Parse()

		if err != io.EOF {
			break
		}
	}
}

func (hlcupRequest *HlcupCtx) Parse() (err error) {
	reader := bufio.NewReader(hlcupRequest.Connection)
	//var body []byte = []byte{}
	body := make([]byte, 1024)

	//reader.Reset(hlcupRequest.Connection)
	_, err = reader.Read(body)

	fmt.Println("START ---------")
	fmt.Println(string(body))
	fmt.Println("END ---------")
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		err = io.EOF
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
		hlcupRequest.PostBody = body[postBodyIndex :]

		//fmt.Println(string(hlcupRequest.PostBody))
	}

	return
}

func (hlcupRequest *HlcupCtx) Error(status int) {
	switch status {
	case 404:
		hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
	case 400:
		hlcupRequest.Connection.Write([]byte("HTTP/1.1 400 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
	default:
		hlcupRequest.Connection.Write([]byte("HTTP/1.1 404 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))
	}

	hlcupRequest.Connection.Close()
}

func (hlcupRequest *HlcupCtx) SetBody(b []byte) {
	hlcupRequest.Connection.Write([]byte("HTTP/1.1 200 OK\ncontent-type:application/json;charset=utf-8;Connection: Closed\n\n"))

	hlcupRequest.Connection.Write(b)
}

func (hlcupRequest *HlcupCtx) Write(b []byte) {
	hlcupRequest.Connection.Write(b)
}

func (hlcupRequest *HlcupCtx) WriteString(s string) {
	hlcupRequest.Connection.Write([]byte(s))
}

func (hlcupRequest *HlcupCtx) Close() {
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