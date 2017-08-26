package logger

import (
	//"bytes"
	"log"
	//"fmt"
	"os"
)

func PrintLog(message string) {

	f, err := os.OpenFile("request.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	log.Println(message)
}