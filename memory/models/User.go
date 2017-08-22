package models

import (
	"time"
	"encoding/json"
	"database/sql/driver"
	"strconv"
	"fmt"
	"github.com/rannoch/highloadcup2017/mysql_server/util"
)

type User struct {
	Id         int32 `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date int32 `json:"birth_date"`
	Visits     []*Visit `json:"-"`
}

func (user *User) HasForeignRelations() bool {
	return false
}

func (user *User) TableName() string {
	return "user"
}

func (user *User) GetId() int32 {
	return user.Id
}

func (user *User) GetFields(alias string) []string {
	return []string{"id", "email", "first_name", "last_name", "gender", "birth_date"}
}

func (user *User) ValidateParams(params map[string]interface{}, scenario string) (result bool) {
	if scenario == "insert" && len(params) != len(user.GetFields("")) {
		return false
	}

	for param, value := range params {
		if value == nil {
			return false
		}

		if scenario == "update" && param == "id" {
			return false
		}

		if !util.StringInSlice(param, user.GetFields("")) {
			return false
		}
	}

	return true
}

func (user *User) SetParams(params map[string]interface{}) {
	id, ok := params["id"]; if ok {
		switch id.(type) {
		case int32:
			user.Id = id.(int32)
		case float32:
			user.Id = int32(id.(float32))
		case float64:
			user.Id = int32(id.(float64))
		case json.Number:
			t, _ := id.(json.Number).Int64()
			user.Id = int32(t)
		}
	}
	email, ok := params["email"]; if ok {
		user.Email = email.(string)
	}
	first_name, ok := params["first_name"]; if ok {
		user.First_name = first_name.(string)
	}
	last_name, ok := params["last_name"]; if ok {
		user.Last_name = last_name.(string)
	}
	gender, ok := params["gender"]; if ok {
		user.Gender = gender.(string)
	}
	birth_date, ok := params["birth_date"]; if ok {
		switch birth_date.(type) {
		case int32:
			user.Birth_date = birth_date.(int32)
		case float32:
			user.Birth_date = int32(birth_date.(float32))
		case float64:
			user.Birth_date = int32(birth_date.(float64))
		case json.Number:
			t, _ := birth_date.(json.Number).Int64()
			user.Birth_date = int32(t)
		}
	}
}

func (user *User) GetValues() []interface{} {
	return []interface{}{user.Id, user.Email, user.First_name, user.Last_name, user.Gender, user.Birth_date}
}

func (user *User) GetFieldPointers(with []string) []interface{} {
	return []interface{}{&user.Id, &user.Email, &user.First_name, &user.Last_name, &user.Gender, &user.Birth_date}
}

type BaskaTime struct {
	time.Time
}

func (baskaTime *BaskaTime) UnmarshalJSON(data []byte) error {
	var timestamp int64
	err := json.Unmarshal(data, &timestamp)
	if err != nil {
		return err
	}

	baskaTime.Time = time.Unix(0, 0)

	baskaTime.Time = baskaTime.Add(time.Duration(timestamp) * time.Second)

	return nil
}

func (baskaTime BaskaTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	//stamp := fmt.Sprintf("\"%s\"", baskaTime.Time.Format(time.RFC1123Z))

	//stamp := baskaTime.Unix() - time.Unix(0,0).Unix()
	return []byte(strconv.Itoa(int(baskaTime.Unix()))), nil
}

func (baskaTime *BaskaTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, ok := value.(time.Time)
	if ok {
		baskaTime.Time = v
	}

	return nil
}

func (baskaTime BaskaTime) Value() (driver.Value, error) {
	return baskaTime.Time, nil
}

type FloatPrecision5 float32

func (f FloatPrecision5) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.5f", f)), nil
}

func (user *User) GetBytes() []byte {
	return []byte(fmt.Sprintf("{\"id\":%d,\"email\":\"%s\",\"first_name\":\"%s\",\"last_name\":\"%s\",\"gender\":\"%s\",\"birth_date\":%d}",
		user.Id, user.Email, user.First_name, user.Last_name, user.Gender, user.Birth_date,
	))
}

func (user *User) GetString() string {
	return fmt.Sprintf("{\"id\":%d,\"email\":\"%s\",\"first_name\":\"%s\",\"last_name\":\"%s\",\"gender\":\"%s\",\"birth_date\":%d}",
		user.Id, user.Email, user.First_name, user.Last_name, user.Gender, user.Birth_date,
	)
}

