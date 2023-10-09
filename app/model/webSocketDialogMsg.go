package model

import (
	"hios/core"
	"hios/utils/common"
	"log"
	"reflect"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
)

// WebSocketDialogMsg 对话消息记录
type WebSocketDialogMsg struct {
	Id         int         `gorm:"primary_key;auto_increment" json:"id"`
	DialogId   int         `gorm:"default:0;comment:对话ID" json:"dialog_id"`
	DialogType string      `gorm:"default:'';comment:对话类型" json:"dialog_type"`
	Userid     int         `gorm:"default:0;comment:发送会员ID" json:"userid"`
	Type       string      `gorm:"default:'';comment:消息类型" json:"type"`
	Mtype      string      `gorm:"default:'';comment:消息类型（用于搜索）" json:"mtype"`
	Msg        string      `gorm:"longtext;comment:详细消息" json:"msg"`
	Emoji      string      `gorm:"type:longtext;comment:emoji回复" json:"emoji"`
	Key        string      `gorm:"type:text;default:'';comment:搜索关键词" json:"key"`
	Read       int         `gorm:"default:0;comment:已阅数量" json:"read"`
	Send       int         `gorm:"default:0;comment:发送数量" json:"send"`
	Tag        int         `gorm:"default:0;comment:标注会员ID" json:"tag"`
	Todo       int         `gorm:"default:0;comment:设为待办会员ID" json:"todo"`
	Link       int         `gorm:"default:0;comment:是否存在链接" json:"link"`
	Modify     int         `gorm:"default:0;comment:是否编辑" json:"modify"`
	ReplyNum   int         `gorm:"default:0;comment:有多少条回复" json:"reply_num"`
	ReplyId    int         `gorm:"default:0;comment:回复ID" json:"reply_id"`
	DeletedAt  core.TsTime `gorm:"comment:删除时间" json:"deleted_at"`
	core.BaseAtModels
	AppendAttrs            map[string]int            `gorm:"-" json:"-"` // 附加属性
	WebSocketDialogMsgRead []*WebSocketDialogMsgRead `gorm:"foreignKey:DialogId" json:"dialog_msg_reads,omitempty"`
	WebSocketDialog        *WebSocketDialog          `gorm:"foreignKey:DialogId" json:"dialog,omitempty"`
}

var WebSocketDialogMsgModel = WebSocketDialogMsg{}

// GetMsgByDialogIdAndMsgId 根据对话ID和消息ID获取对话消息
func (m WebSocketDialogMsg) GetMsgByDialogIdAndMsgId(dialogID, msgID int) (*WebSocketDialogMsg, error) {
	err := core.DB.Where("dialog_id = ? AND id = ?", dialogID, msgID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetLatestMsgByDialogId 根据对话ID获取最新的对话消息
func (m WebSocketDialogMsg) GetLatestMsgByDialogId(dialogID int) (*WebSocketDialogMsg, error) {
	err := core.DB.Where("dialog_id = ?", dialogID).Order("id DESC").First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// WebSocketDialog 返回关联WebSocketDialog
// func (m WebSocketDialogMsg) WebSocketDialog() (*WebSocketDialog, error) {
// 	dialog := &WebSocketDialog{}
// 	err := core.DB.Where("id = ?", m.DialogId).First(dialog).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dialog, nil
// }

// GetPercentage 阅读占比
func (m WebSocketDialogMsg) GetPercentage() int {
	if _, ok := m.AppendAttrs["percentage"]; !ok {
		m.GeneratePercentage()
	}
	return m.AppendAttrs["percentage"]
}

// Increment 自增
func (m WebSocketDialogMsg) Increment(column string, amount int) {
	err := core.DB.Model(&m).UpdateColumn(column, core.Expr(column+" + ?", amount)).Error
	if err != nil {
		panic(err.Error())
	}
}

// GeneratePercentage 生成阅读占比
func (m *WebSocketDialogMsg) GeneratePercentage(increment ...int) int {
	log.Printf("GeneratePercentage: %v", m)
	if len(increment) > 0 {
		add := increment[0]
		m.Increment("read", add)
	}
	if m.AppendAttrs == nil {
		m.AppendAttrs = make(map[string]int)
	}
	if m.Read > m.Send || m.Send == 0 {
		m.AppendAttrs["percentage"] = 100
	} else {
		m.AppendAttrs["percentage"] = int(float64(m.Read) / float64(m.Send) * 100)
		m.AppendAttrs["percentage"] = 90
	}
	log.Printf("GeneratePercentage: %v %v", m.Read, m.Send)
	return 100
}

// 前置钩子函数
// func (m *WebSocketDialogMsg) BeforeFind(tx *gorm.DB) (err error) {
// 	// 处理消息数据
// 	m.Msg = fmt.Sprintf("%v", m.GetMsgAttribute())
// 	// 处理emoji数据
// 	m.Emoji = fmt.Sprintf("%v", m.GetEmojiAttribute())
// 	return nil
// }

// 后置钩子函数
// func (m *WebSocketDialogMsg) AfterFind(tx *gorm.DB) (err error) {
// 	// 处理消息数据
// 	m.Msg = fmt.Sprintf("%v", m.GetMsgAttribute())
// 	// 处理emoji数据
// 	m.Emoji = fmt.Sprintf("%v", m.GetEmojiAttribute())
// 	return nil
// }

// 消息格式化
func (m WebSocketDialogMsg) GetMsgAttribute() interface{} {
	var value interface{}
	if m.Msg != "" {
		value, _ = common.StrToMap(m.Msg)
		if m.Type == "file" {
			if v, ok := value.(map[string]interface{}); ok {
				ext, _ := v["ext"].(string)
				v["type"] = "file"
				if common.InArray(ext, []string{"jpg", "jpeg", "webp", "png", "gif"}) {
					v["type"] = "img"
				}
				v["path"] = common.FillUrl(v["path"].(string))
				if thumb, ok := v["thumb"].(string); ok {
					v["thumb"] = common.FillUrl(thumb)
				} else {
					v["thumb"] = common.FillUrl(common.ExtIcon(ext))
				}
			}
		} else if m.Type == "record" {
			if v, ok := value.(map[string]interface{}); ok {
				v["path"] = common.FillUrl(v["path"].(string))
			}
		}
	}
	return value
}

// emoji回复格式化
func (m WebSocketDialogMsg) GetEmojiAttribute() interface{} {
	if common.IsKind(m.Msg, reflect.Map) {
		return m.Msg
	}
	value, _ := common.StrToMap(m.Msg)
	return value
}

// GenerateMsgKey 生成关键词
func (m *WebSocketDialogMsg) GenerateMsgKey() string {
	switch m.Type {
	case "text":
		value, _ := common.StrToMap(m.Msg)
		return strings.ReplaceAll(strip.StripTags(value["text"].(string)), "&nbsp;", " ")
	case "meeting", "file":
		// todo 获取msg["name"]
		return string(m.Msg)
	default:
		return ""
	}
}

// Update 更新实例
func (m *WebSocketDialogMsg) Update(data map[string]interface{}) error {
	if err := common.MapToStruct(data, &m); err != nil {
		return err
	}
	return core.DB.Updates(&m).Error
}

// Save 保存实例
func (m *WebSocketDialogMsg) Save() error {
	return core.DB.Save(&m).Error
}

// Create 创建实例
func (m WebSocketDialogMsg) Create() error {
	return core.DB.Create(&m).Error
}
