package models

import (
	"time"
	"encoding/json"
	"database/sql/driver"
	"strconv"
)

type User struct {
	Id         int32 `json:"id"`
	Email      string `json:"email"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Gender     string `json:"gender"`
	Birth_date BaskaTime `json:"birth_date"`
}

func (user *User) TableName() string{
	return "user"
}

func (user *User) GetId() int32 {
	return user.Id
}

func (user *User) GetFields() []string{
	return []string{"id", "email", "first_name", "last_name", "gender", "birth_date"}
}

func (user *User) GetValues() []interface{} {
	return []interface{}{user.Id, user.Email, user.First_name, user.Last_name, user.Gender, user.Birth_date}
}

func (user *User) GetFieldPointers() []interface{} {
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

	baskaTime.Time = time.Unix(0,0)

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
	v, ok := value.(time.Time); if ok {
		baskaTime.Time = v
	}

	return nil
}


func (baskaTime BaskaTime) Value() (driver.Value, error) {
	return baskaTime.Time, nil
}