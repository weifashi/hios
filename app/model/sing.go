package model

import "hios/core"

// 签名表
type Sing struct {
	Id   int    `gorm:"primary_key;auto_increment" json:"id"`
	Sing string `gorm:"type:varchar(36);default:'';comment:签名" json:"sing"`
	Use  bool   `gorm:"type:tinyint(1);default:0;comment:是否使用" json:"use"`
	Md5  string `gorm:"type:varchar(36);default:'';comment:md5" json:"md5"`
	core.BaseAtModels
}

var SingModel = Sing{}
