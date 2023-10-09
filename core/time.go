package core

import (
	"encoding/json"
	"time"
)

const DateFormat = "2006-01-02"
const TimeFormat = "2006-01-02 15:04:05"

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
