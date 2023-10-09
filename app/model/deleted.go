package model

import (
	"hios/core"
)

type Deleted struct {
	core.BaseIdModels
	Type   string `gorm:"type:varchar(50);default:;comment:删除的数据类型（如：project、task、dialog）" json:"type"`
	Did    int    `gorm:"default:0;comment:删除的数据ID" json:"did"`
	Userid int    `gorm:"default:0;comment:成员ID" json:"userid"`
	core.BaseAtModels
}

var DeletedModel = Deleted{}

// 获取Ids
func (m Deleted) GetIds(types string, userid int, times ...any) []int {
	var ids []int
	db := core.DB.Model(&m).Where("type = ?", types).Where("userid = ?", userid).Order("id desc")
	//
	if len(times) == 0 || times[0] == 0 {
		db = db.Limit(50)
	} else {
		db = db.Where("created_at >= ?", times[0]).Limit(500)
	}
	db.Pluck("did", &ids)
	//
	return ids
}
