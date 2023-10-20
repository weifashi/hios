package interfaces

import "github.com/gorilla/websocket"

type WebsocketReq struct {
	Token    string `json:"token"`
	Language string `json:"language"`
	Type     string `json:"type"`
}

type WsClient struct {
	Conn *websocket.Conn `json:"conn"`

	Ip   string `json:"ip"`   // 客户端IP
	Type string `json:"type"` // 客户端类型（如：user）
	Uid  string `json:"uid"`  // 客户端用户ID（会员ID）
	Rid  int32  `json:"rid"`  // 客户端序号ID（WebSocket ID）
}

type WsMsg struct {
	Action int `json:"action"` // 消息类型：1、上线；2、下线；3、消息
	Data   any `json:"data"`   // 消息内容

	Type string `json:"type"` // 客户端类型（如：user）
	Uid  string `json:"uid"`  // 客户端用户ID（会员ID）
	Rid  int32  `json:"rid"`  // 客户端序号ID（WebSocket ID）

	Md5    string `json:"md5"`
	Output string `json:"output"`
	Err    string `json:"err"`
}
