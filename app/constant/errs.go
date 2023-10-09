package constant

// api
var (
	// 环境
	ErrEnvProhibition   = "ErrEnvProhibition"   //当前环境禁止此操作
	ErrInvalidParameter = "ErrInvalidParameter" //参数错误

	//用户
	ErrTypeNotLogin           = "ErrTypeNotLogin"           //未登录
	ErrNoPermission           = "ErrNoPermission"           //无权限
	ErrNotSupport             = "ErrNotSupport"             //不支持的请求方式
	ErrConnFailed             = "ErrConnectionFailed"       //连接失败
	ErrPasswordFailed         = "ErrPasswordFailed"         //账号或密码错误
	ErrPasswordEnterFailed    = "ErrPasswordEnterFailed"    //请输入登录密码
	ErrCaptchaCode            = "ErrCaptchaCode"            //验证码错误
	ErrRegFail                = "ErrRegFail"                //注册失败
	ErrRegSuccess             = "ErrRegSuccess"             //注册成功，请验证邮箱后登录
	ErrNotOpenReg             = "ErrNotOpenReg"             //未开放注册
	ErrInvitationCode         = "ErrInvitationCode"         //请输入正确的邀请码
	ErrEmailAddress           = "ErrEmailAddress"           //请输入正确的邮箱地址
	ErrRegVerificationMailbox = "ErrRegVerificationMailbox" //您的帐号已注册过，请验证邮箱
	ErrEmailAddressExists     = "ErrEmailAddressExists"     //邮箱地址已存在
	ErrEmailAddressError      = "ErrEmailAddressError"      //邮箱地址错误
	ErrEmailInconformity      = "ErrEmailInconformity"      //与当前登录邮箱不一致
	ErrInputNewEmailAddress   = "ErrInputNewEmailAddress"   //请输入新邮箱地址
	ErrOldEmailAddressAccord  = "ErrOldEmailAddressAccord"  //不能与旧邮箱一致

	// 用户部门
	ErrDeptsNameLength        = "ErrDeptNameLength"        //部门名称长度限制2-20个字
	ErrDeptsNameSpecialSymbol = "ErrDeptNameSpecialSymbol" //部门名称不能包含特殊符号
	ErrDeptsNameM             = "ErrDeptsNameM"            //部门名称不能包含：(M)
	ErrDeptsInexistence       = "ErrDeptsInexistence"      //部门不存在或已被删除
	ErrDeptsMaxRestrict       = "ErrDeptsMaxRestrict"      //最多只能创建200个部门
	ErrUpDeptsInexistence     = "ErrUpDeptsInexistence"    //上级部门不存在或已被删除
	ErrUpDeptsHierarchyError  = "ErrUpDeptsHierarchyError" //上级部门层级错误
	ErrMaxSubDeptRestrict     = "ErrMaxSubDeptRestrict"    //每个部门最多只能创建20个子部门
	ErrExistSubDeptNotEdit    = "ErrExistSubDeptNotEdit"   //含有子部门无法修改上级部门
	ErrDepartmentHead         = "ErrDepartmentHead"        //请选择正确的部门负责人

	// 系统设置
	ErrEnvForbidden           = "ErrEnvForbidden"           //当前环境禁止修改
	ErrAutoArchiveTime        = "ErrAutoArchiveTime"        //自动归档时间不可小于1天！
	ErrAutoArchiveTimeMax     = "ErrAutoArchiveTimeMax"     //自动归档时间不可大于100天！
	ErrAutoArchiveTimeInvalid = "ErrAutoArchiveTimeInvalid" //自动归档时间不可小于1天或大于100天！
	ErrRequestDataEmpty       = "ErrRequestDataEmpty"       //请求数据不能为空
	ErrRequestDataInvalid     = "ErrRequestDataInvalid"     //请求数据不合法
	ErrMailToEmpty            = "ErrMailToEmpty"            //收件人地址不能为空
	ErrMailHostEmpty          = "ErrMailHostEmpty"          //SMTP服务器不能为空
	ErrMailPortEmpty          = "ErrMailPortEmpty"          //SMTP端口不能为空
	ErrMailUserEmpty          = "ErrMailUserEmpty"          //SMTP账号不能为空
	ErrMailPassEmpty          = "ErrMailPassEmpty"          //SMTP密码不能为空
	ErrMailToInvalid          = "ErrMailToInvalid"          //请输入正确的收件人地址
	ErrMailNotConfig          = "ErrMailNotConfig"          //发送邮箱未配置
	ErrRequestTimeout         = "ErrRequestTimeout"         //请求超时
	ErrMailContentReject      = "ErrMailContentReject"      //邮件内容被拒绝，请检查邮箱是否开启接收功能
	ErrUserEmailEmpty         = "ErrUserEmailEmpty"         //用户邮箱不能为空
	ErrUserTelEmpty           = "ErrUserTelEmpty"           //用户联系电话不能为空
	ErrUserNameEmpty          = "ErrUserNameEmpty"          //用户昵称不能为空
	ErrUserIdEmpty            = "ErrUserIdEmpty"            //用户ID不能为空
	ErrUserOldPasswordEmpty   = "ErrUserOldPasswordEmpty"   //用户旧密码不能为空
	ErrUserNewPasswordEmpty   = "ErrUserNewPasswordEmpty"   //用户新密码不能为空
	ErrUserPasswordSame       = "ErrUserPasswordSame"       //新旧密码不能一致
	ErrUserOldPassword        = "ErrUserOldPassword"        //请填写正确的旧密码

	// 聊天
	ErrRecordNotFound      = "ErrRecordNotFound"      //记录不存在
	ErrChatSessionNotExist = "ErrChatSessionNotExist" //会话不存在或已被删除
	ErrMyMessage           = "ErrMyMessage"           //@我的消息
	ErrSetNickname         = "ErrSetNickname"         //请设置昵称
	ErrSetTel              = "ErrSetTel"              //请设置联系电话
)

// app
var (
	ErrCmdTimeout = "ErrCmdTimeout" //命令执行超时
)
