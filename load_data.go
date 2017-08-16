package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/rannoch/highloadcup2017/models"
	"encoding/json"
	"log"
	"github.com/rannoch/highloadcup2017/storage"
)

func LoadData(path string) (err error) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		fmt.Println(f.Name())

		entities := parseFile(path + "/" + f.Name())

		/*for _, entity := range entities {
			err = storage.InsertEntity(entity)

			if err != nil {
				log.Println(err.Error())
			}
		}*/
		err = storage.InsertEntityMultiple(entities)
		if err != nil {
			log.Println(err.Error())
		}
	}

	return
}

func parseFile(file string) (entities []models.Entity) {
	var err error

	fileContent, err := ioutil.ReadFile(file)

	if err != nil {
		log.Println(err.Error())
		return
	}

	switch {
	case strings.Contains(file, "users"):
		var m = make(map[string][]models.User)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			log.Println(err.Error())
			return
		}
		for _, v := range m["users"] {
			entities = append(entities, v)
		}

	case strings.Contains(file, "locations"):
		var m = make(map[string][]models.Location)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			log.Println(err.Error())
			return
		}

		for _, v := range m["locations"] {
			entities = append(entities, v)
		}
	case strings.Contains(file, "visits"):
		var m = make(map[string][]models.Visit)
		err = json.Unmarshal(fileContent, &m)

		if err != nil {
			log.Println(err.Error())
			return
		}

		for _, v := range m["visits"] {
			entities = append(entities, v)
		}
	}

	return
}
