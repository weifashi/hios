package system

import (
	"hios/app/constant"
	"hios/app/interfaces"
	"hios/app/model"
	"hios/utils/common"
	"hios/utils/email"
	e "hios/utils/error"
	"net/mail"

	"os"
)

var SystemService = systemService{}

type systemService struct{}

// SystemSetting 系统设置
func (s systemService) SystemSetting(param interfaces.SystemSettingReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	paramType := param.Type
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		// 校验自动归档时间
		if param.AutoArchived == "open" && (param.ArchivedDay <= 0 || param.ArchivedDay > 100) {
			return nil, e.New(constant.ErrAutoArchiveTimeInvalid) // 自动归档时间不可小于1天或大于100天！
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation(model.SettingSystemKey, settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation(model.SettingSystemKey, nil, false)
		ret = data
	}

	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"reg":              "open",
			"reg_identity":     "normal",
			"login_code":       "auto",
			"password_policy":  "simple",
			"project_invite":   "open",
			"chat_information": "optional",
			"anon_message":     "open",
			"auto_archived":    "close",
			"archived_day":     7,
			"all_group_mute":   "open",
			"all_group_autoin": "yes",
			"image_compress":   "open",
			"image_save_local": "open",
			"start_home":       "close",
		}
	}
	//
	if regInvite, ok := ret["reg_invite"].(string); !ok || regInvite == "" {
		if paramType == "all" || paramType == "save" {
			ret["reg_invite"] = common.GeneratePassword(8, "")
		} else {
			ret["reg_invite"] = ""
		}
	}
	return ret, nil
}

// SystemSettingMail 系统邮件设置
func (s systemService) SystemSettingMail(param interfaces.SystemSettingMailReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation(model.SettingEmailKey, settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation(model.SettingEmailKey, nil, false)
		ret = data
	}
	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"smtp_server":             "",
			"port":                    "",
			"account":                 "",
			"password":                "",
			"reg_verify":              "close",
			"notice_msg":              "close",
			"msg_unread_user_minute":  -1,
			"msg_unread_group_minute": -1,
			"ignore_addr":             "",
		}
	}
	return ret, nil
}

// SystemSettingMeeting 系统会议设置
func (s systemService) SystemSettingMeeting(param interfaces.SystemSettingMeetingReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation("meetingSetting", settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("meetingSetting", nil, false)
		ret = data
	}
	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"open": "close",
		}

	} else {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			// todo appid、app_certificate安全显示
			return nil, nil
		}
	}
	return ret, nil
}

// SystemSettingCheckin 系统签到设置
func (s systemService) SystemSettingCheckin(param interfaces.SystemSettingCheckinReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		if param.Open == "close" {
			param.Key = common.StringMd5(common.GeneratePassword(32, ""))
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation("checkinSetting", settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("checkinSetting", nil, false)
		ret = data
	}
	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"open":         "close",
			"time":         "00:00",
			"advance":      120,
			"delay":        120,
			"remindin":     5,
			"remindexceed": 10,
			"edit":         "close",
		}
	}
	ret["cmd"] = "curl -sSL url | sh"
	return ret, nil
}

// SystemSettingApppush 系统应用推送设置
func (s systemService) SystemSettingApppush(param interfaces.SystemSettingApppushReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation("apppushSetting", settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("apppushSetting", nil, false)
		ret = data
	}
	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"push": "close",
		}
	}
	return ret, nil
}

// SystemSettingThirdaccess 系统第三方访问设置
func (s systemService) SystemSettingThirdaccess(param interfaces.SystemSettingThirdaccessReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	// todo 测试第三方访问
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		// 保存设置
		settingMap, _ := common.StructToMap(param)
		delete(settingMap, "type")
		data, err := model.SettingModel.SettingOperation("thirdaccessSetting", settingMap, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("thirdaccessSetting", nil, false)
		ret = data
	}
	// 如果ret为空，初始化数据
	if ret == nil {
		ret = map[string]interface{}{
			"ldap_open":       "close",
			"ldap_port":       389,
			"ldap_sync_local": "close",
		}
	}
	return ret, nil
}

// SystemSettingPriority 系统优先级设置
func (s systemService) SystemSettingPriority(param interfaces.SystemSettingPriorityReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		//
		if len(param.List) == 0 {
			return nil, e.New(constant.ErrRequestDataEmpty) //请求数据不能为空
		}
		var listMapSlice []map[string]interface{}
		for _, item := range param.List {
			if item.Name == "" || item.Color == "" || item.Days <= 0 || item.Priority <= 0 {
				return nil, e.New(constant.ErrRequestDataInvalid) //请求数据不合法
			}
			listMapSlice = append(listMapSlice, map[string]interface{}{
				"name":     item.Name,
				"color":    item.Color,
				"days":     item.Days,
				"priority": item.Priority,
			})
		}
		// 保存设置
		data, err := model.SettingModel.SettingOperation("prioritySetting", map[string]interface{}{"list": listMapSlice}, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("prioritySetting", nil, false)
		ret = data
	}
	return ret, nil
}

// SystemSettingTemplate 系统模板设置
func (s systemService) SystemSettingTemplate(param interfaces.SystemSettingTemplateReq) (map[string]interface{}, error) {
	// todo 检查用户权限
	var ret map[string]interface{}
	//
	if param.Type == "save" {
		if env := os.Getenv("SYSTEM_SETTING"); env == "disabled" {
			return nil, e.New(constant.ErrEnvForbidden) //当前环境禁止修改
		}
		//
		if len(param.List) == 0 {
			return nil, e.New(constant.ErrRequestDataEmpty) //请求数据不能为空
		}
		var listMapSlice []map[string]interface{}
		for _, item := range param.List {
			if item.Name == "" || item.Columns == "" {
				return nil, e.New(constant.ErrRequestDataInvalid) //请求数据不合法
			}
			listMapSlice = append(listMapSlice, map[string]interface{}{
				"name":    item.Name,
				"columns": item.Columns,
			})
		}
		// 保存设置
		data, err := model.SettingModel.SettingOperation("templateSetting", map[string]interface{}{"list": listMapSlice}, true)
		if err != nil {
			return nil, err
		}
		ret = data

	} else {
		// 获取设置
		data, _ := model.SettingModel.SettingOperation("templateSetting", nil, false)
		ret = data
	}
	return ret, nil
}

// SystemEmailCheck 用于测试配置邮箱是否能发送邮件
func (s systemService) SystemEmailCheck(to string, smtpServer string, port int, account string, password string) error {

	subject := "Mail sending test"
	body := `<p>收到此电子邮件意味着您的邮箱配置正确。</p><p>Receiving this email means that your mailbox is configured correctly.</p>`

	// 验证收件人地址是否合法
	if _, err := mail.ParseAddress(to); err != nil {
		return e.New(constant.ErrMailToInvalid)
	}

	email := email.NewEmailService(smtpServer, port, account, password)

	// 发送信息
	err := email.Send(to, subject, body)
	if err != nil {
		return err
	}

	return nil
}
