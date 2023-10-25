package interfaces

type SeedReq struct {
	Uid    string `json:"uid" binding:"required"` // rid
	Path   string `json:"path" binding:""`        // 执行文件路径
	Source string `json:"source" binding:""`      // 类型
	Type   string `json:"type" binding:""`        // 类型
	Force  bool   `json:"force" binding:""`       // 是否强制发送（默认5秒内不能发送重覆指令）
	Msg    string `json:"msg" binding:"required"` // 消息内容
	Before string `json:"before" binding:""`      // 执行前先执行的内容
	After  string `json:"after" binding:""`       // 执行后执行的内容
}
