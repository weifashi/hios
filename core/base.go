package core

import (
	"encoding/json"
	"time"
)

const DateFormat = "2006-01-02"
const TimeFormat = "2006-01-02 15:04:05"

type BaseIdModels struct {
	Id int `gorm:"primary_key" json:"id"`
}

type BaseAtModels struct {
	CreatedAt *time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TsTime 自定义时间格式
type TsTime int64

func (tst *TsTime) UnmarshalJSON(bs []byte) error {
	var date string
	err := json.Unmarshal(bs, &date)
	if err != nil {
		return err
	}
	if date == "" || date == "null" {
		return nil
	}
	tt, _ := time.ParseInLocation(TimeFormat, date, time.Local)
	*tst = TsTime(tt.Unix())
	return nil
}

func (tst TsTime) MarshalJSON() ([]byte, error) {
	if tst == 0 {
		return json.Marshal(nil)
	}
	tt := time.Unix(int64(tst), 0).Format(TimeFormat)
	return json.Marshal(tt)
}
