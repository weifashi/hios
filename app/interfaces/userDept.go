package interfaces

// 部门列表
type UserDeptInfo struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`         //身份证
	DialogID    int    `json:"dialog_id"`    //对话ID
	ParentId    int    `json:"parent_id"`    //上级部门ID
	OwnerUserId int    `json:"owner_userid"` //部门负责人ID
	CreatedAt   int64  `json:"created_at"`   //创建时间
	UpdatedAt   int64  `json:"updated_at"`   //更新时间
}

// 添加部门
type UserDeptAddReq struct {
	Id          int    `json:"id"`           //部门id，留空为创建部门
	Name        string `json:"name"`         //部门名称
	ParentId    int    `json:"parent_id"`    //上级部门ID
	OwnerUserId int    `json:"owner_userid"` //部门负责人ID
	DialogGroup string `json:"dialog_group"` //部门群（仅创建部门时有效） new: 创建（默认） use: 使用现有群
	DialogUseid int    `json:"dialog_useid"` //使用现有群ID（dialog_group=use时有效）
}
