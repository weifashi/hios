package interfaces

type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"` //验证码ID
	ImagePath string `json:"image_path"` //图片路径
}

type UserInfo struct {
	Id        int    `json:"id"`         // 用户id
	Email     string `json:"email"`      // 邮箱
	Name      string `json:"name"`       // 用户名称
	Token     string `json:"token"`      // Token
	Avatar    string `json:"avatar"`     // 头像
	CreatedAt string `json:"created_at"` // 创建时间
	UpdatedAt string `json:"updated_at"` // 更新时间
}

type UserLoginReq struct {
	Email    string `json:"email" binding:"required"`     //邮箱
	Password string `json:"password" binding:"required" ` //密码
	CodeId   string `json:"code_id"`                      //验证码ID
	Code     string `json:"code"`                         //验证码
}

type UserRegReq struct {
	Email    string `json:"email" binding:"required"`    //邮箱
	Password string `json:"password" binding:"required"` //密码
	Invite   string `json:"invite"`                      //邀请码
}

type UserLoginQrcodeReq struct {
	Code string `json:"code" binding:"required"` //二维码 code
}

// UserEditDataReq 个人设置
type UserEditDataReq struct {
	Avatar     string `json:"avatar"`                   // 头像
	Email      string `json:"email" binding:"required"` // 邮箱
	Tel        string `json:"tel" binding:"required"`   // 联系电话
	Name       string `json:"name" binding:"required"`  // 用户名称
	Profession string `json:"profession"`               // 职位/职称
}

// UserEditPasswordReq 密码设置
type UserEditPasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// UserEditEmailReq 修改邮箱
type UserEditEmailReq struct {
	Email string `json:"email" binding:"required"` // 邮箱
}

// UserDeleteAccount 删除账号
type UserDeleteAccountReq struct {
	Email    string `json:"email" binding:"required"`    // 帐号邮箱
	Code     string `json:"code" binding:"required"`     // 邮箱验证码
	Reason   string `json:"reason" binding:"required"`   // 注销理由
	Password string `json:"password" binding:"required"` // 注销理由
	Type     string `json:"type" binding:"required"`     // 类型
}

// UserSearch 搜索
type UserSearchReq struct {
	*Pages
	Keys        map[string]interface{} `json:"keys"`         // 搜索条件: keys.key 昵称、拼音、邮箱关键字;  keys.disable 0-排除离职（默认），1-仅离职，2-含离职 ; keys.bot 0-排除机器人（默认），1-仅机器人，2-含机器人; keys.project_id 在指定项目ID; keys.no_project_id 不在指定项目ID; keys.dialog_id 在指定对话ID
	Sorts       map[string]interface{} `json:"sorts"`        // 排序方式: sorts.az 按字母：asc|desc
	UpdatedTime int64                  `json:"updated_time"` // 在这个时间戳之后更新的
	State       int                    `json:"state"`        // 获取在线状态 0: 不获取（默认）,1: 获取会员在线状态，返回数据多一个online值
}

// UserListReq 会员列表
type UserListReq struct {
	*Pages
	Keys          map[string]interface{} `json:"keys"`            // 搜索条件: keys.key 邮箱/电话/昵称/职位; keys.email 邮箱; keys.tel 电话; keys.name 昵称; keys.profession 职位; keys.identity 身份（如：admin、noadmin）;  keys.disable 是否离职 yes: 仅离职 all:全部 其他值: 仅在职（默认); keys.email_verity 邮箱是否认证 yes no; keys.bot 是否包含机器人 yes:仅机器人 all: 全部 其他值: 非机器人（默认）;
	GetCheckinMac int                    `json:"get_checkin_mac"` // 获取签到mac地址  0: 不获取（默认）1: 获取
}

type UserOperationMacReq struct {
	Mac    string `json:"mac"`    //mac地址
	Remark string `json:"remark"` //备注
}

// UserOperationReq 操作会员（限管理员）
type UserOperationReq struct {
	Userid         int                   `json:"userid"`          // 会员ID
	Type           string                `json:"type"`            // 操作 : setadmin 设为管理员, clearadmin 取消管理员, settemp 设为临时帐号,cleartemp 取消临时身份（取消临时帐号）,checkin_macs  修改自动签到mac地址需要参数checkin_macs, department 修改部门需要参数department, setdisable 设为离职（需要参数 disable_time、transfer_userid）, cleardisable 取消离职, delete 删除会员（需要参数 delete_reason)
	Email          string                `json:"email"`           //邮箱地址
	Tel            string                `json:"tel"`             //联系电话
	Password       string                `json:"password"`        //新的密码
	Name           string                `json:"name"`            //昵称
	Profession     string                `json:"profession"`      //职位
	CheckinMacs    []UserOperationMacReq `json:"checkin_macs"`    //自动签到mac地址
	Department     []int                 `json:"department"`      //部门
	DisableTtime   string                `json:"disable_time"`    //部门
	TransferUserid int                   `json:"transfer_userid"` //离职交接人
	DeleteReason   string                `json:"delete_reason"`   //删除原因
}
