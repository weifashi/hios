package interfaces

type SeedReq struct {
	Email    string `json:"email" binding:"required"`    // 帐号邮箱
	Code     string `json:"code" binding:"required"`     // 邮箱验证码
	Reason   string `json:"reason" binding:"required"`   // 注销理由
	Password string `json:"password" binding:"required"` // 注销理由
	Type     string `json:"type" binding:"required"`     // 类型
}
