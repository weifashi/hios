package interfaces

import (
	"hios/app/model"
	"hios/core"
)

// DialogListReq 对话列表请求参数
type DialogListReq struct {
	Updated int64 `json:"updated"` // 更新时间
	Deleted int64 `json:"deleted"` // 删除时间
	*Pages
}

// DialogSearchReq 对话搜索请求参数
type DialogSearchReq struct {
	Key string `json:"key"` // 关键词
}

// DialogMsgSendTextReq 对话消息发送请求参数
type DialogMsgSendTextReq struct {
	DialogId int    `json:"dialog_id" binding:"required"` // 会话ID
	Text     string `json:"text" binding:"required"`      // 消息内容
	TextType string `json:"text_type"`                    // 消息类型 html：HTML文本 md：Markdown文本
	UpdateId int    `json:"update_id"`                    // 更新消息ID（优先大于 reply_id）
	ReplyId  int    `json:"reply_id"`                     // 回复ID
	Silence  string `json:"silence"`                      // 是否静默发送 no：否 yes：是
}

// DialogMsgSendVoiceReq 对话消息发送语音请求参数
type DialogMsgSendVoiceReq struct {
	DialogId int    `json:"dialog_id" binding:"required"` // 会话ID
	ReplyId  int    `json:"reply_id"`                     // 回复ID
	Base64   string `json:"base64" binding:"required"`    // 语音base64编码
	Duration int    `json:"duration" binding:"required"`  // 语音时长（毫秒）
}

// DialogMsgSendFileReq 对话消息发送文件请求参数
type DialogMsgSendFileReq struct {
	DialogId        int    `form:"dialog_id" binding:"required"` // 会话ID
	ReplyId         int    `form:"reply_id"`                     // 回复ID
	ImageAttachment int    `form:"image_attachment"`             // 图片是否也存到附件
	Filename        string `form:"filename"`                     // post-文件名称
	Base64          string `form:"base64"`                       // post-base64（二选一）
}

// DialogMsgSendFilesReq 对话消息发送多个文件请求参数
type DialogMsgSendFilesReq struct {
	Userids         string `form:"userids"`          // 发送给的成员ID，格式: userid1, userid2, userid3 （dialog_ids 二选一）
	DialogIds       string `form:"dialog_ids"`       // 对话ID（user_ids 二选一）
	ReplyId         int    `form:"reply_id"`         // 回复ID
	ImageAttachment int    `form:"image_attachment"` // 图片是否也存到附件
	Filename        string `form:"filename"`         // post-文件名称
	Base64          string `form:"base64"`           // post-base64（二选一）
}

// DialogMsgSendFileIdReq 对话消息发送文件ID请求参数
type DialogMsgSendFileIdReq struct {
	FileId    int   `json:"file_id" binding:"required"`    // 文件ID
	DialogIds []int `json:"dialog_ids" binding:"required"` // 转发给的对话ID
	Userids   []int `json:"userids" binding:"required"`    // 转发给的成员ID
}

// DialogMsgListReq 对话消息列表请求参数
type DialogMsgListReq struct {
	DialogId   int    `json:"dialog_id"`                   // 会话ID
	MsgId      int    `json:"msg_id"`                      // 消息ID
	PositionId int    `json:"position_id"`                 // 此消息ID前后的数据
	PrevId     int    `json:"prev_id"`                     // 此消息ID之前的数据
	NextId     int    `json:"next_id"`                     // 此消息ID之后的数据  - position_id、prev_id、next_id 只有一个有效，优先循序为：position_id > prev_id > next_id
	MsgType    string `json:"msg_type"`                    // 消息类型 tag: 标记 link: 链接 text: 文本 image: 图片 file: 文件 record: 录音 meeting: 会议
	Take       int    `form:"take,default=50" json:"take"` // 获取条数，默认:50，最大:100
}

// DialogMsgSearchReq 对话消息搜索请求参数
type DialogMsgSearchReq struct {
	DialogId int    `json:"dialog_id"` // 会话ID
	Key      string `json:"key"`       // 搜索关键词
}

// DialogMsgForwardReq 对话消息转发请求参数
type DialogMsgForwardReq struct {
	MsgId     int   `json:"msg_id" binding:"required"`     // 消息ID
	DialogIds []int `json:"dialog_ids" binding:"required"` // 转发给的对话ID
	Userids   []int `json:"userids" binding:"required"`    // 转发给的成员ID
}

// DialogGroupAddReq 对话群组添加请求参数
type DialogGroupAddReq struct {
	Avatar   string `json:"avatar"`                     // 群头像
	ChatName string `json:"chat_name"`                  // 群名称
	Userids  []int  `json:"userids" binding:"required"` // 群成员，格式: [userid1, userid2, userid3]
}

// DialogGroupEditReq 对话群组编辑请求参数
type DialogGroupEditReq struct {
	DialogId int    `json:"dialog_id" binding:"required"` // 会话ID
	Avatar   string `json:"avatar"`                       // 群头像
	ChatName string `json:"chat_name"`                    // 群名称
	Admin    int    `json:"admin"`                        // 系统管理员操作（1：只判断是不是系统管理员，否则判断是否群管理员）
}

// FormatDialog 对话格式化结构
type FormatDialog struct {
	*model.WebSocketDialog                            // WebSocketDialog 模型的指针
	Email                  string                     `json:"email"`                   // 对话中的电子邮件地址
	UserAt                 core.TsTime                `json:"user_at"`                 // 对话中的用户最后一次活动时间
	UserMs                 core.TsTime                `json:"user_ms"`                 // 对话中的用户最后一次活动时间戳（毫秒）
	TopAt                  core.TsTime                `json:"top_at"`                  // 对话中的置顶时间戳
	SearchMsgId            *int                       `json:"search_msg_id,omitempty"` // 对话中要搜索的消息 ID
	LastMsg                *model.WebSocketDialogMsg  `json:"last_msg,omitempty"`      // 对话中的最后一条消息
	LastAt                 core.TsTime                `json:"last_at"`                 // 对话中的最后一条消息的时间戳（毫秒）
	MarkUnread             int                        `json:"mark_unread"`             // 是否标记为未读：0否，1是
	Silence                int                        `json:"silence"`                 // 是否免打扰：0否，1是
	People                 int64                      `json:"people"`                  // 对话中的参与人数
	TodoNum                int64                      `json:"todo_num"`                // 对话中的待办事项数量
	Pinyin                 string                     `json:"pinyin"`                  // 对话名称的拼音
	QuickMsgs              []map[string]string        `json:"quick_msgs"`              // 对话中的快捷消息
	DialogUser             *model.WebSocketDialogUser `json:"dialog_user,omitempty"`   // 对话中的用户信息
	GroupInfo              interface{}                `json:"group_info,omitempty"`    // 对话中的群组信息
	Bot                    int                        `json:"bot"`                     // 对话中是否存在机器人
	DialogDelete           int                        `json:"dialog_delete"`           // 对话中是否已删除
	AllGroupMute           string                     `json:"all_group_mute"`          // 对话中是否已开启全员禁言
	HasTag                 bool                       `json:"has_tag"`                 // 对话中是否存在标签消息
	HasImage               bool                       `json:"has_image"`               // 对话中是否存在图片消息
	HasFile                bool                       `json:"has_file"`                // 对话中是否存在文件消息
	HasLink                bool                       `json:"has_link"`                // 对话中是否存在链接消息
	Unread                 int64                      `json:"unread"`                  // 对话中未读消息的数量
	Mention                int64                      `json:"mention"`                 // 对话中@我的消息的数量
	PositionMsgs           interface{}                `json:"position_msgs,omitempty"` // 对话中的位置消息
}
