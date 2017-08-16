package main

import (
	"io/ioutil"
	"fmt"
)

func load_data(path string) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())
	}
}