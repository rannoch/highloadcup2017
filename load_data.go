package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/rannoch/highloadcup2017/models"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/storage"
	"time"
	"sync"
)

func LoadData(path string) (err error) {
	start := time.Now()

	files, _ := ioutil.ReadDir(path)
	var wg sync.WaitGroup

	for _, f := range files {
		entities := parseFile(path + "/" + f.Name())

		fmt.Printf("%s len - %d \n", f.Name(), len(entities))

		if len(entities) > 5000 {
			for i := 0; i < len(entities); i = i + 5000 {
				start := i
				end := i + 5000
				if end > len(entities) {
					end = len(entities)
				}

				go func() {
					wg.Add(1)
					storage.Db.InsertEntityMultiple(entities[start: end])
					wg.Done()
				}()
			}
		} else {
			go func() {
				wg.Add(1)
				storage.Db.InsertEntityMultiple(entities)
				wg.Done()
			} ()
		}
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println(elapsed)

	return
}

func parseFile(file string) (entities []storage.Entity) {
	var err error

	fileContent, err := ioutil.ReadFile(file)

	if err != nil {
		return
	}

	switch {
	case strings.Contains(file, "users"):
		var m = make(map[string][]models.User)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}
		for _, v := range m["users"] {
			c := v
			entities = append(entities, &c)
		}

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.Location)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["locations"] {
			c := v
			entities = append(entities, &c)
		}
	case strings.Contains(file, "visits"):
		var m = make(map[string][]models.Visit)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["visits"] {
			c := v
			entities = append(entities, &c)
		}
	}

	return
}
