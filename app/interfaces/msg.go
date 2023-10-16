package interfaces

type SeedReq struct {
	Rid    string `json:"rid" binding:"required"` // rid
	Path   string `json:"path" binding:""`        // 执行文件路径
	Cmd    string `json:"cmd" binding:"required"` // 执行内容
	Before string `json:"before" binding:""`      // 执行前先执行的内容
	After  string `json:"after" binding:""`       // 执行后执行的内容
}
