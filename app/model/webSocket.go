package model

import "hios/core"

// WebSocket WebSocket连接
type WebSocket struct {
	Id   int    `gorm:"primary_key;auto_increment" json:"id"`
	Key  string `gorm:"not null;unique_index;comment:WebSocket连接的唯一标识" json:"key"`
	Fd   string `gorm:"default:'';comment:WebSocket连接的文件描述符" json:"fd"`
	Path string `gorm:"default:'';comment:WebSocket连接的路径" json:"path"`
	Uid  string `gorm:"default:'';comment:WebSocket连接的用户ID,节点ip" json:"uid"`
	core.BaseAtModels
}

var WebSocketModel = WebSocket{}

// UpdateInsert 更新或插入
func (m WebSocket) UpdateInsert(where, data map[string]interface{}) error {
	return core.DB.Table(core.DBTableName(&WebSocket{})).Where(where).Assign(data).FirstOrCreate(&WebSocket{}).Error
}
