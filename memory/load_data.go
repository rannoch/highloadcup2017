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
	for _, visit := range storage.VisitDb {
		user, ok := storage.UserDb[visit.User]

		if ok {
			visit.User_model = user

			user.Visits = append(user.Visits, visit)
		}

		location, ok := storage.LocationDb[visit.Location]

		if ok {
			visit.Location_model = location

			location.Visits = append(location.Visits, visit)
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

			storage.UserDb[v.Id] = &c
		}

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.Location)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["locations"] {
			c := v

			storage.LocationDb[v.Id] = &c
		}
	case strings.Contains(file, "visits"):
		var m = make(map[string][]models.Visit)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["visits"] {
			c := v

			storage.VisitDb[v.Id] = &c
		}
	}

	return
}
