package model

import "hios/core"

// WebSocketTmpMsg WebSocket临时消息
type WebSocketTmpMsg struct {
	Id       int    `gorm:"primary_key;auto_increment" json:"id"`
	Md5      string `gorm:"default:'';comment:MD5(会员ID-消息)" json:"md5"`
	Msg      string `gorm:"type:longtext;comment:详细消息" json:"msg"`
	Send     int    `gorm:"default:0;comment:是否已发送" json:"send"`
	CreateId string `gorm:"default:'';comment:所属会员ID" json:"create_id"`
	core.BaseAtModels
}

var WebSocketTmpMsgModel = WebSocketTmpMsg{}
