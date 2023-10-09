package model

import (
	"errors"
	"fmt"
	"hios/core"
)

// WebSocketDialog WebSocket对话
type WebSocketDialog struct {
	Id        int         `gorm:"primary_key;auto_increment" json:"id"`
	Type      string      `gorm:"default:'';comment:对话类型" json:"type"`
	GroupType string      `gorm:"default:'';comment:聊天室类型" json:"group_type"`
	Name      string      `gorm:"default:'';comment:对话名称" json:"name"`
	Avatar    string      `gorm:"default:'';comment:头像（群）" json:"avatar"`
	OwnerId   int         `gorm:"default:0;comment:群主用户ID" json:"owner_id"`
	LastAt    core.TsTime `gorm:"comment:最后消息时间" json:"last_at"`
	DeletedAt core.TsTime `gorm:"comment:删除时间" json:"deleted_at"`
	core.BaseAtModels
	WebSocketDialogUser []*WebSocketDialogUser `gorm:"foreignKey:DialogId" json:"dialog_users,omitempty"`
}

var WebSocketDialogModel = WebSocketDialog{}

// GetAvatar 返回对话的头像URL
func (m *WebSocketDialog) GetAvatar() string {
	if m.Avatar != "" {
		return BaseFillURL(m.Avatar)
	}
	return m.Avatar
}

// BaseFillURL 填充给定URL的基本URL
// todo 实现填充基本URL的逻辑
func BaseFillURL(url string) string {
	return url
}

// DialogUser 返回对话用户
func (m WebSocketDialog) DialogUser() []*WebSocketDialogUser {
	var users []*WebSocketDialogUser
	core.DB.Where("dialog_id = ?", m.Id).Find(&users)
	return users
}

// GetWebSocketDialogByID 根据ID获取对话
func (m WebSocketDialog) GetWebSocketDialogByID(id int) (*WebSocketDialog, error) {
	err := core.DB.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetWebSocketDialogByIds 根据用户ID获取对话
func (m WebSocketDialog) GetWebSocketDialogByIds(userid, targetUserId int) (*WebSocketDialog, error) {
	var dialog WebSocketDialog

	webSocketDialogTable := core.DBTableName(&WebSocketDialog{})
	webSocketDialogUserTable := core.DBTableName(&WebSocketDialogUser{})

	err := core.DB.Table(webSocketDialogTable+" AS w").
		Where("type = ?", "user").
		Joins(fmt.Sprintf("JOIN %s AS u1 ON w.id = u1.dialog_id", webSocketDialogUserTable)).
		Joins(fmt.Sprintf("JOIN %s AS u2 ON w.id = u2.dialog_id", webSocketDialogUserTable)).
		Where("u1.userid = ? AND u2.userid = ?", userid, targetUserId).
		First(&dialog).Error
	if err != nil {
		if errors.Is(err, core.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dialog, nil
}

// CreateWebSocketDialog 创建对话
func (m WebSocketDialog) CreateWebSocketDialog(userid, targetUserId int) (*WebSocketDialog, error) {
	dialog := &WebSocketDialog{
		Type: "user",
	}
	err := core.DB.Create(dialog).Error
	if err != nil {
		return nil, err
	}
	err = core.DB.Create(&WebSocketDialogUser{
		DialogId: dialog.Id,
		Userid:   userid,
	}).Error
	if err != nil {
		return nil, err
	}
	err = core.DB.Create(&WebSocketDialogUser{
		DialogId: dialog.Id,
		Userid:   targetUserId,
	}).Error
	if err != nil {
		return nil, err
	}
	return dialog, nil
}

// CheckMute 检查禁言
func (m WebSocketDialog) CheckMute(userid int) {
	if m.GroupType == "all" {
		system, _ := SettingModel.GetSetting(SettingSystemKey)
		allGroupMute := system["all_group_mute"].(string)
		switch allGroupMute {
		case "all":
			panic("当前会话全员禁言")
		case "user":
			var user User
			if err := core.DB.Model(&User{}).Where("id = ?", userid).First(&user).Error; err != nil {
				panic(err)
			}
			if !user.IsAdmin() {
				panic("当前会话禁言")
			}
		}
	}
}

// UpdateMsgLastAt 更新对话最后消息时间
func (m *WebSocketDialog) UpdateMsgLastAt() *WebSocketDialogMsg {
	lastMsg := &WebSocketDialogMsg{}
	if err := core.DB.Where("dialog_id = ?", m.Id).Order("id desc").First(lastMsg).Error; err != nil {
		return nil
	}

	m.LastAt = lastMsg.CreatedAt
	if err := core.DB.Save(m).Error; err != nil {
		return nil
	}

	return lastMsg
}
