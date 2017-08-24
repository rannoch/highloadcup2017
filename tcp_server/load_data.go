package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/rannoch/highloadcup2017/tcp_server/models"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/tcp_server/storage"
	"time"
	"sync"
	"sort"
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
		if visit == nil {
			continue
		}

		user := storage.UserDb[visit.User]

		visit.User_model = user

		user.Visits = append(user.Visits, visit)

		// сортировка заранее
		sort.Sort(models.VisitByDateAsc(user.Visits))

		location := storage.LocationDb[visit.Location]

		visit.Location_model = location

		location.Visits = append(location.Visits, visit)
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
			storage.UserBytesDb[v.Id] = c.GetBytes()
		}
		storage.UserCount += int64(len(m["users"]))

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.Location)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["locations"] {
			c := v

			storage.LocationDb[v.Id] = &c
			storage.LocationBytesDb[v.Id] = c.GetBytes()
		}
		storage.LocationCount += int64(len(m["locations"]))

	case strings.Contains(file, "visits"):
		var m = make(map[string][]models.Visit)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["visits"] {
			c := v

			storage.VisitDb[v.Id] = &c
			storage.VisitBytesDb[v.Id] = c.GetBytes()
		}

		storage.VisitCount += int64(len(m["visits"]))
	}

	return
}
