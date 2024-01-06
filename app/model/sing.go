package model

import (
	"hios/core"
	"time"
)

// 签名表
type Sing struct {
	Id    int        `gorm:"primary_key;auto_increment" json:"id"`
	Sing  string     `gorm:"type:varchar(36);default:'';comment:签名" json:"sing"`
	UseAt *time.Time `gorm:"comment:使用时间" json:"use_at"`
	Md5   string     `gorm:"type:varchar(36);default:'';comment:md5" json:"md5"`
	core.BaseAtModels
}

var SingModel = Sing{}
