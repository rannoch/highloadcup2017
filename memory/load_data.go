package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/rannoch/highloadcup2017/memory/models"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/memory/storage"
	"time"
	"sync"
)

func LoadData(path string) (err error) {
	start := time.Now()

	files, _ := ioutil.ReadDir(path)
	var wg sync.WaitGroup

	for _, f := range files {
		parseAndAppendFile(path + "/" + f.Name())
	}

	// dependencies
	for _, v := range storage.Db["visit"] {
		visit := v.(*models.Visit)

		user, ok := storage.Db["user"][visit.User_id]

		if ok {
			visit.User = user.(*models.User)

			user.(*models.User).Visits = append(user.(*models.User).Visits, visit)
		}

		location, ok := storage.Db["location"][visit.Location_id]

		if ok {
			visit.Location = location.(*models.Location)

			location.(*models.Location).Visits = append(location.(*models.Location).Visits, visit)
		}
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println(elapsed)

	return
}

func parseAndAppendFile(file string) () {
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

			storage.Db["user"][v.Id] = &c
		}

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.Location)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["locations"] {
			c := v

			storage.Db["location"][v.Id] = &c
		}
	case strings.Contains(file, "visits"):
		var m = make(map[string][]models.Visit)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["visits"] {
			c := v

			storage.Db["visit"][v.Id] = &c
		}
	}

	return
}
