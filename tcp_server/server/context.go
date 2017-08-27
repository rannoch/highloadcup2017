package server

import (
	"net"
	"github.com/valyala/fasthttp"
	"bufio"
	"bytes"
	"strconv"
	"fmt"
	"github.com/rannoch/highloadcup2017/tcp_server/logger"
	"io"
)

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

	ResponseBodyBuffer *bytes.Buffer
	ResponseFullBuffer *bytes.Buffer

	Server *TcpServer
}

func NewCtx(server *TcpServer, connection net.Conn) *HlcupCtx {
	responseBuffer := server.ResponseBodyBufferPool.Get().(*bytes.Buffer)
	responseBuffer.Reset()

	responseFullBuffer := server.ResponseFullBufferPool.Get().(*bytes.Buffer)
	responseFullBuffer.Reset()

	return &HlcupCtx{
		Connection:         connection,
		ResponseBodyBuffer: responseBuffer,
		ResponseFullBuffer: responseFullBuffer,
		ResponseStatus:     200,
		Server: server,
	}
}

func (hlcupRequest *HlcupCtx) ResetParams() {
	hlcupRequest.ResponseStatus = 200

	hlcupRequest.ResponseBodyBuffer.Reset()
	hlcupRequest.ResponseFullBuffer.Reset()
	hlcupRequest.QueryArgs.Reset()
}

func (hlcupRequest *HlcupCtx) Handle() {
	//net.TCPConn.SetKeepAlive(true)

	//reader := bufio.NewReader(hlcupRequest.Connection)

	reader := hlcupRequest.Server.ReaderPool.Get().(*bufio.Reader)

	for {
		buf := hlcupRequest.Server.BytePool.Get().([]byte)

		reader.Reset(hlcupRequest.Connection)
		n, err := reader.Read(buf)

		hlcupRequest.Server.BytePool.Put(buf)

		//logger.PrintLog(string(buf))
		//time.Sleep(10 * time.Millisecond)
		//fmt.Println(hlcupRequest.ConnId)

		if err != nil && err != io.EOF {
			hlcupRequest.Server.ReaderPool.Put(reader)
			fmt.Println(err.Error())
			logger.PrintLog(fmt.Sprintf("%d %s %s", hlcupRequest.ConnId, err.Error(), string(buf)))
			hlcupRequest.Close()
			hlcupRequest.Connection = nil
			return
		}

		hlcupRequest.ResetParams()
		err = hlcupRequest.Parse(buf, n)

		hlcupRequest.Server.HandleFunc(hlcupRequest)

		if hlcupRequest.ConnectionClosed {
			hlcupRequest.Connection = nil
			hlcupRequest.Server.ReaderPool.Put(reader)
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
	writer := hlcupRequest.Server.WriterPool.Get().(*bufio.Writer)
	writer.Reset(hlcupRequest.Connection)

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


			writer.Write(hlcupRequest.ResponseFullBuffer.Bytes())
			writer.Flush()
			hlcupRequest.Server.WriterPool.Put(writer)

			//hlcupRequest.Connection.Write(hlcupRequest.ResponseFullBuffer.Bytes())
			hlcupRequest.Server.ResponseFullBufferPool.Put(hlcupRequest.ResponseFullBuffer)
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


			writer.Write(hlcupRequest.ResponseFullBuffer.Bytes())
			writer.Flush()
			hlcupRequest.Server.WriterPool.Put(writer)

			hlcupRequest.Server.ResponseFullBufferPool.Put(hlcupRequest.ResponseFullBuffer)

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

	hlcupRequest.ResponseFullBuffer.WriteString(strconv.Itoa(len(hlcupRequest.ResponseBodyBuffer.Bytes())))
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)
	hlcupRequest.ResponseFullBuffer.Write(strCRLF)

	hlcupRequest.ResponseFullBuffer.Write(hlcupRequest.ResponseBodyBuffer.Bytes())

	hlcupRequest.Server.ResponseBodyBufferPool.Put(hlcupRequest.ResponseBodyBuffer)

	writer.Write(hlcupRequest.ResponseFullBuffer.Bytes())
	writer.Flush()
	hlcupRequest.Server.WriterPool.Put(writer)

	//logger.PrintLog(fmt.Sprintf("- %d - %s", hlcupRequest.ConnId, hlcupRequest.ResponseFullBuffer.String()))

	hlcupRequest.Server.ResponseFullBufferPool.Put(hlcupRequest.ResponseFullBuffer)

	if !hlcupRequest.KeepAlive {
		hlcupRequest.Close()
	}
}
