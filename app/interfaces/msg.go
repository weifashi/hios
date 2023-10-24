package interfaces

type SeedReq struct {
	Uid    string `json:"uid" binding:"required"` // rid
	Path   string `json:"path" binding:""`        // 执行文件路径
	Type   string `json:"type" binding:""`        // 类型
	Msg    string `json:"msg" binding:"required"` // 消息内容
	Before string `json:"before" binding:""`      // 执行前先执行的内容
	After  string `json:"after" binding:""`       // 执行后执行的内容
}
