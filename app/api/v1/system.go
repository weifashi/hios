package v1

import (
	"hios/app/constant"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/app/service/system"
	"strconv"
)

// @Tags System
// @Summary 获取/保存系统设置
// @Description 获取/保存系统设置
// @Accept json
// @Param request body interfaces.SystemSettingReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingReq}
// @Router /system/setting [post]
func (api *BaseApi) SystemSetting() {
	var param = interfaces.SystemSettingReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSetting(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存邮箱设置（限管理员）
// @Description 获取/保存邮箱设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingMailReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingMailReq}
// @Router /system/setting/mail [post]
func (api *BaseApi) SystemSettingMail() {
	var param = interfaces.SystemSettingMailReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingMail(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存会议设置（限管理员）
// @Description 获取/保存会议设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingMeetingReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingMeetingReq}
// @Router /system/setting/meeting [post]
func (api *BaseApi) SystemSettingMeeting() {
	var param = interfaces.SystemSettingMeetingReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingMeeting(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存签到设置（限管理员）
// @Description 获取/保存签到设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingCheckinReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingCheckinReq}
// @Router /system/setting/checkin [post]
func (api *BaseApi) SystemSettingCheckin() {
	var param = interfaces.SystemSettingCheckinReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingCheckin(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存应用推送设置（限管理员）
// @Description 获取/保存应用推送设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingApppushReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingApppushReq}
// @Router /system/setting/apppush [post]
func (api *BaseApi) SystemSettingApppush() {
	var param = interfaces.SystemSettingApppushReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingApppush(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存第三方账号设置（限管理员）
// @Description 获取/保存第三方账号设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingThirdaccessReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingThirdaccessReq}
// @Router /system/setting/thirdaccess [post]
func (api *BaseApi) SystemSettingThirdaccess() {
	var param = interfaces.SystemSettingThirdaccessReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingThirdaccess(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存任务优先级设置（限管理员）
// @Description 获取/保存任务优先级设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingPriorityReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingPriorityReq}
// @Router /system/setting/priority [post]
func (api *BaseApi) SystemSettingPriority() {
	var param = interfaces.SystemSettingPriorityReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingPriority(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 获取/保存保存项目模板设置（限管理员）
// @Description 获取/保存保存项目模板设置（限管理员）
// @Accept json
// @Param request body interfaces.SystemSettingTemplateReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingTemplateReq}
// @Router /system/setting/template [post]
func (api *BaseApi) SystemSettingTemplate() {
	var param = interfaces.SystemSettingTemplateReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	result, err := system.SystemService.SystemSettingTemplate(param)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, result)
}

// @Tags System
// @Summary 邮箱设置检查
// @Description 邮箱设置检查
// @Accept json
// @Param request body interfaces.SystemSettingMailToReq true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingMailToReq}
// @Router /system/email/check [post]
func (api *BaseApi) SystemEmailCheck() {
	var param = interfaces.SystemSettingMailToReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	// 判断值是否为空
	if param.To == "" {
		helper.ErrorWith(api.Context, constant.ErrMailToEmpty, nil)
		return
	}
	if param.SmtpServer == "" {
		helper.ErrorWith(api.Context, constant.ErrMailHostEmpty, nil)
		return
	}
	if strconv.Itoa(param.Port) == "" {
		helper.ErrorWith(api.Context, constant.ErrMailPortEmpty, nil)
		return
	}
	if param.Account == "" {
		helper.ErrorWith(api.Context, constant.ErrMailUserEmpty, nil)
		return
	}
	if param.Password == "" {
		helper.ErrorWith(api.Context, constant.ErrMailPassEmpty, nil)
		return
	}

	err := system.SystemService.SystemEmailCheck(param.To, param.SmtpServer, param.Port, param.Account, param.Password)
	if err != nil {
		helper.ErrorWith(api.Context, err.Error(), nil)
		return
	}
	helper.Success(api.Context, "")
}
