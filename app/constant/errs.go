package constant

// 环境
var (
	ErrEnvProhibition    = "当前环境禁止此操作"
	ErrInvalidParameter  = "参数错误"
	ErrNotSupport        = "不支持的请求方式"
	ErrConnFailed        = "连接失败"
	ErrRequestTimeout    = "请求超时"
	ErrMailContentReject = "邮件内容被拒绝，请检查邮箱是否开启接收功能"
	ErrMailNotConfig     = "发送邮箱未配置"
	ErrCaptchaCode       = "验证码错误"
	ErrTypeNotLogin      = "未登录"
	ErrNoPermission      = "无权限"
	ErrCmdTimeout        = "命令执行超时"
)
