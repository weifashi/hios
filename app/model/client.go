package model

import "hios/core"

// WebSocket WebSocket连接
type Client struct {
	Id     int    `gorm:"primary_key;auto_increment" json:"id"`
	Uid    string `gorm:"type:varchar(36);default:'';comment:客户端标识" json:"uid"`
	Source string `gorm:"type:varchar(50);default:'';comment:来源" json:"source"`
	IP     string `gorm:"type:varchar(50);default:'';comment:客户端ip" json:"ip"`
	Online bool   `gorm:"type:tinyint(1);default:0;comment:是否在线" json:"online"`
	Swap   int    `gorm:"default:0;comment:客户端虚拟内存 (MB)" json:"swap"`
	Cc     int    `gorm:"default:0;comment:累计连接次数" json:"cc"`
	core.BaseAtModels
}

var ClientModel = Client{}

// UpdateInsert 更新或插入
func (m Client) UpdateInsert(where, data map[string]interface{}) error {
	return core.DB.Table(core.DBTableName(&Client{})).Where(where).Assign(data).FirstOrCreate(&Client{}).Error
}
