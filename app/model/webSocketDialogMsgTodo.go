package model

import "hios/core"

// WebSocketDialogMsgTodo 对话消息待办记录
type WebSocketDialogMsgTodo struct {
	Id       int         `gorm:"primary_key;auto_increment" json:"id"`
	DialogId int         `gorm:"default:0;comment:对话ID" json:"dialog_id"`
	MsgId    int         `gorm:"default:0;comment:消息ID" json:"msg_id"`
	Userid   int         `gorm:"default:0;comment:接收会员ID" json:"userid"`
	DoneAt   core.TsTime `gorm:"comment:完成时间" json:"done_at"`
	core.BaseAtModels
}

var WebSocketDialogMsgTodoModel = WebSocketDialogMsgTodo{}

// GetMsgTodoCountByDialogIdAndUserId 根据对话ID和用户ID获取待办消息数量
func (m WebSocketDialogMsgTodo) GetMsgTodoCountByDialogIdAndUserId(dialogId, userid int) (int64, error) {
	var count int64
	err := core.DB.Model(&m).Where("dialog_id = ? AND userid = ? AND (done_at IS NULL OR done_at = 0)", dialogId, userid).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
