package model

import (
	"hios/core"
)

// WebSocketDialogUser 对话用户
type WebSocketDialogUser struct {
	Id         int         `gorm:"primary_key;auto_increment" json:"id"`
	DialogId   int         `gorm:"default:0;comment:对话ID" json:"dialog_id"`
	Userid     int         `gorm:"default:0;comment:会员ID" json:"userid"`
	MarkUnread int         `gorm:"default:0;comment:是否标记为未读：0否，1是" json:"mark_unread"`
	Silence    int         `gorm:"default:0;comment:是否免打扰：0否，1是" json:"silence"`
	Inviter    int         `gorm:"default:0;comment:邀请人" json:"inviter"`
	Important  int         `gorm:"default:0;comment:是否不可移出（项目、任务、部门人员）" json:"important"`
	TopAt      core.TsTime `gorm:"comment:置顶时间" json:"top_at"`
	core.BaseAtModels
	WebSocketDialog *WebSocketDialog `gorm:"foreignKey:DialogId" json:"dialog,omitempty"`
}

var WebSocketDialogUserModel = WebSocketDialogUser{}

// CheckDialogUser 检查对话用户
func (m WebSocketDialogUser) CheckDialogUser(dialogId, userid int) bool {
	var count int64
	core.DB.Model(&m).
		Where("dialog_id = ? AND userid = ?", dialogId, userid).
		Count(&count)
	return count > 0
}

// GetDialogUser 获取对话用户 by userID
func (m WebSocketDialogUser) GetDialogUser(userid int) (*WebSocketDialogUser, error) {
	err := core.DB.Model(&m).
		Where("userid = ?", userid).
		First(&m).Error
	return &m, err
}

// GetByDialogIdAndUserId 根据对话ID和用户ID获取对话用户
func (m WebSocketDialogUser) GetByDialogIdAndUserId(dialogId, userid int) (*WebSocketDialogUser, error) {
	err := core.DB.Where("dialog_id = ? AND userid = ?", dialogId, userid).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetCountByDialogId 根据对话ID获取对话用户数量
func (m WebSocketDialogUser) GetCountByDialogId(dialogId int) (int64, error) {
	var count int64
	err := core.DB.Model(&WebSocketDialogUser{}).Where("dialog_id = ?", dialogId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetByDialogIdAndNotUserId 根据对话ID和用户ID获取不是该用户的对话用户
func (m WebSocketDialogUser) GetByDialogIdAndNotUserId(dialogId, userid int) (*WebSocketDialogUser, error) {
	err := core.DB.Where("dialog_id = ? AND userid != ?", dialogId, userid).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
