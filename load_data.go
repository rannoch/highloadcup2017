package main

import (
	"io/ioutil"
	"fmt"
)

func LoadData(path string) {


	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())


		dat, err := ioutil.ReadFile("/tmp/dat")
	}
}

func parseFile(file string) {

}