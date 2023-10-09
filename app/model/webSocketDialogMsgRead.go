package model

import "hios/core"

// WebSocketDialogMsgRead 对话消息阅读记录
type WebSocketDialogMsgRead struct {
	Id       int         `gorm:"primary_key;auto_increment" json:"id"`
	DialogId int         `gorm:"default:0;comment:对话ID" json:"dialog_id"`
	MsgId    int         `gorm:"default:0;comment:消息ID" json:"msg_id"`
	Userid   int         `gorm:"default:0;comment:接收会员ID" json:"userid"`
	Mention  int         `gorm:"default:0;comment:是否提及（被@）" json:"mention"`
	Silence  int         `gorm:"default:0;comment:是否免打扰：0否，1是" json:"silence"`
	Email    int         `gorm:"default:0;comment:是否发了邮件" json:"email"`
	After    int         `gorm:"default:0;comment:在阅读之后才添加的记录" json:"after"`
	ReadAt   core.TsTime `gorm:"comment:阅读时间" json:"read_at"`
	core.BaseAtModels
	WebSocketDialogMsg *WebSocketDialogMsg `gorm:"foreignKey:MsgId" json:"dialog_msgs,omitempty"`
}

var WebSocketDialogMsgReadModel = WebSocketDialogMsgRead{}
