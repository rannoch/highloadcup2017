package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/rannoch/highloadcup2017/epoll_server/models"
	"encoding/json"
	"github.com/rannoch/highloadcup2017/epoll_server/storage"
	"time"
	"sync"
	"sort"
	"github.com/antonholmquist/jason"
	"strconv"
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
	case strings.Contains(file, "options"):
		optionsStr := string(fileContent)

		generateTime, err := strconv.ParseInt(strings.Split(optionsStr, "\n")[0], 0, 64)

		if err != nil {
			fmt.Println(err.Error())
			storage.GenerateTime = time.Now()
			break
		}

		storage.GenerateTime = time.Unix(generateTime, 0)
	case strings.Contains(file, "users"):
		v, err := jason.NewObjectFromBytes(fileContent)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		users, err := v.GetObjectArray("users")

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, object := range users {
			email, _ := object.GetString("email")
			first_name, _ := object.GetString("first_name")
			last_name, _ := object.GetString("last_name")
			gender, _ := object.GetString("gender")

			c := models.User{}
			c.Id, _ = object.GetInt64("id")
			c.Email = []byte(email)
			c.First_name = []byte(first_name)
			c.Last_name = []byte(last_name)
			c.Gender = []byte(gender)
			c.Birth_date, _ = object.GetInt64("birth_date")

			storage.UserDb[c.Id] = &c
			//storage.UserBytesDb[c.Id] = c.GetBytes()
		}
		storage.UserCount += int64(len(users))

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.LocationHelper)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			return
		}

		for _, v := range m["locations"] {
			c := v.GetLocation()

			storage.LocationDb[v.Id] = &c
			//storage.LocationBytesDb[v.Id] = c.GetBytes()
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
			//storage.VisitBytesDb[v.Id] = c.GetBytes()
		}

		storage.VisitCount += int64(len(m["visits"]))
	}

	return
}
