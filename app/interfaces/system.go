package interfaces

// SystemSettingReq 系统设置
type SystemSettingReq struct {
	Type            string  `json:"type"`                  //类型 get: 获取 (默认)，save: 保存，all: 获取所有
	Reg             string  `json:"reg"`                   //注册设置 open: 开启 (默认)，close: 关闭，invite: 邀请
	RegIdentity     string  `json:"reg_identity"`          //注册身份 normal：正常 (默认)，temp：临时
	RegInvite       string  `json:"reg_invite,omitempty"`  //注册邀请码
	LoginCode       string  `json:"login_code"`            //登录验证码设置 auto: 自动 (默认)，open: 开启，close: 关闭，
	PasswordPolicy  string  `json:"password_policy"`       //密码策略设置 simple: 简单 (默认)，complex: 复杂
	ProjectInvite   string  `json:"project_invite"`        //项目邀请设置 open: 开启 (默认)，close: 关闭
	ChatInformation string  `json:"chat_information"`      //聊天资料设置 optional：可选 (默认)，required：必填
	AnonMessage     string  `json:"anon_message"`          //匿名消息设置 open: 开启 (默认)，close: 关闭
	AutoArchived    string  `json:"auto_archived"`         //自动归档设置 open: 开启 (默认)，close: 关闭
	ArchivedDay     float64 `json:"archived_day"`          //自动归档天数
	AllGroupMute    string  `json:"all_group_mute"`        //全员禁言设置 open: 开放 (默认)，user: 成员禁言，all：全员禁言
	AllGroupAutoin  string  `json:"all_group_autoin"`      //自动进入全员群 yes: 是 (默认)，no: 否
	ImageCompress   string  `json:"image_compress"`        //图片优化设置 open: 开启 (默认)，close: 关闭
	ImageSaveLocal  string  `json:"image_save_local"`      //保存网络图片设置 open: 开启 (默认)，close: 关闭
	StartHome       string  `json:"start_home"`            //是否启动首页 open: 开启 (默认)，close: 关闭
	HomeFooter      string  `json:"home_footer,omitempty"` //首页底部设置
}

// SystemSettingMailReq 系统邮件设置
type SystemSettingMailReq struct {
	Type                 string `json:"type"`                    //类型 get: 获取 (默认)，save: 保存
	RegVerify            string `json:"reg_verify"`              //开启注册验证设置 open: 开启 (默认)，close: 关闭
	NoticeMsg            string `json:"notice_msg"`              //消息提醒设置 open: 开启 (默认)，close: 关闭
	MsgUnreadUserMinute  int    `json:"msg_unread_user_minute"`  //用户消息未读提醒时间
	MsgUnreadGroupMinute int    `json:"msg_unread_group_minute"` //群组消息未读提醒时间
	IgnoreAddr           string `json:"ignore_addr"`             //忽略邮箱
	*SystemSettingMailServerReq
}

// SystemSettingMailServerReq 邮箱服务器设置
type SystemSettingMailServerReq struct {
	SmtpServer string `json:"smtp_server"` //SMTP服务器
	Port       int    `json:"port"`        //端口
	Account    string `json:"account"`     //账号
	Password   string `json:"password"`    //密码
}

// SystemSettingMailToReq 邮箱服务器发送设置
type SystemSettingMailToReq struct {
	To string `json:"to"` //接收邮箱
	*SystemSettingMailServerReq
}

// SystemSettingMeetingReq 系统会议设置
type SystemSettingMeetingReq struct {
	Type           string `json:"type"`            //类型 get: 获取 (默认)，save: 保存
	Open           string `json:"open"`            //开启会议设置 open: 开启，close: 关闭 (默认)
	Appid          string `json:"appid"`           //appid
	AppCertificate string `json:"app_certificate"` //app证书
}

// SystemSettingCheckinReq 系统签到设置
type SystemSettingCheckinReq struct {
	Type         string `json:"type"`         //类型 get: 获取 (默认)，save: 保存
	Open         string `json:"open"`         //开启签到设置 open: 开启，close: 关闭 (默认)
	Time         string `json:"time"`         //签到时间
	Advance      int    `json:"advance"`      //最早可提前
	Delay        int    `json:"delay"`        //最晚可延后
	Remindin     int    `json:"remindin"`     //签到打卡提醒
	Remindexceed int    `json:"remindexceed"` //签到缺卡提醒
	Edit         string `json:"edit"`         //允许修改 open: 开启 (默认)，close: 关闭
	Key          string `json:"key"`          //签到安装口令
}

// SystemSettingApppushReq 系统应用推送设置
type SystemSettingApppushReq struct {
	Type          string `json:"type"`           //类型 get: 获取 (默认)，save: 保存
	Push          string `json:"push"`           //推送设置 open: 开启，close: 关闭 (默认)
	IosKey        string `json:"ios_key"`        //ios key
	IosSecret     string `json:"ios_secret"`     //ios secret
	AndroidKey    string `json:"android_key"`    //android key
	AndroidSecret string `json:"android_secret"` //android secret
}

// SystemSettingThirdaccessReq 系统第三方访问设置
type SystemSettingThirdaccessReq struct {
	Type       string `json:"type"`            //类型 get: 获取 (默认)，save: 保存
	LdapOpen   string `json:"ldap_open"`       //ldap开启设置 open: 开启，close: 关闭 (默认)
	LdapHost   string `json:"ldap_host"`       //ldap地址
	LdapPort   int    `json:"ldap_port"`       //ldap端口
	LdapPasswd string `json:"ldap_password"`   //ldap密码
	LdapUserDn string `json:"ldap_user_dn"`    //ldap用户dn
	LdapBaseDn string `json:"ldap_base_dn"`    //ldap基础dn
	LdapSync   string `json:"ldap_sync_local"` //ldap同步本地用户设置 open: 开启，close: 关闭 (默认)
}

// SystemSettingPriorityReq 系统优先级设置
type SystemSettingPriorityReq struct {
	Type string                       `json:"type"`           //类型 get: 获取 (默认)，save: 保存
	List []*SystemSettingPriorityList `json:"list,omitempty"` //列表
}

// SystemSettingPriorityList 系统优先级设置列表
type SystemSettingPriorityList struct {
	Name     string `json:"name"`     //名称
	Color    string `json:"color"`    //颜色
	Days     int    `json:"days"`     //天数
	Priority int    `json:"priority"` //级别
}

// SystemSettingTemplateReq 系统模板设置
type SystemSettingTemplateReq struct {
	Type string                       `json:"type"`           //类型 get: 获取 (默认)，save: 保存
	List []*SystemSettingTemplateList `json:"list,omitempty"` //列表
}

// SystemSettingTemplateList 系统模板设置列表
type SystemSettingTemplateList struct {
	Name    string `json:"name"`    //名称
	Columns string `json:"columns"` //项目模板
}
