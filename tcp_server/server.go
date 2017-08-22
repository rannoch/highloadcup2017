package main

import "net"
import "fmt"
import (
	"bufio"
	"net/textproto"
	"log"
	"net/http"
)



func main() {

	fmt.Println("Launching server...")

	// listen on all interfaces
	listener, _ := net.Listen("tcp", ":8086")

	// run loop forever (or until ctrl-c)
	for {
		// accept connection on port
		connection, _ := listener.Accept()



		// will listen for message to process ending in newline (\n)
		message := bufio.NewReader(connection)
		// output message received
		// sample process for string received

		tp := textproto.NewReader(message)

		mimeHeader, err := tp.ReadMIMEHeader()
		if err != nil {
			log.Fatal(err)
		}

		// http.Header and textproto.MIMEHeader are both just a map[string][]string
		httpHeader := http.Header(mimeHeader)

		// send new string back to client
		connection.Write([]byte(httpHeader.Get("Content-Type")))
		connection.Close()
	}
}