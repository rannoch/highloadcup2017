package models

import (
	"time"
	"encoding/json"
	"database/sql/driver"
	"strconv"
	"fmt"
	"github.com/rannoch/highloadcup2017/util"
	"reflect"
	"strings"
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
	userValue := reflect.ValueOf(user).Elem()

	for param, value := range params {
		field := userValue.FieldByName(strings.Title(param))
		fmt.Sprintf("%s=%v", param, value)

		switch field.Interface().(type) {
		case int32:
			switch value.(type) {
			case int32:
				field.Set(reflect.ValueOf(value.(int32)))
			case float32:
				field.Set(reflect.ValueOf(int32(value.(float32))))
			case float64:
				field.Set(reflect.ValueOf(int32(value.(float64))))
			}
		case string:
			field.SetString(value.(string))
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
	v, ok := value.(time.Time);
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
