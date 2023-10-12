package v1

import (
	"hios/app/helper"
	"hios/core"
	"strconv"
)

// @Tags System
// @Summary 发送消息
// @Description 发送消息
// @Accept json
// @Param request rid content true "request"
// @Param request body content true "request"
// @Success 200 {object} interfaces.Response{data=interfaces.SystemSettingReq}
// @Router /msg/seed [post]
func (api *BaseApi) MsgSeed() {
	rid := api.Context.PostForm("rid")
	path := api.Context.PostForm("path")
	data := api.Context.PostForm("content")
	types := api.Context.DefaultPostForm("type", "file")
	before := api.Context.DefaultPostForm("before", "")
	after := api.Context.DefaultPostForm("after", "")

	if len(rid) > 0 && len(data) > 0 {
		rd, _ := strconv.Atoi(rid)
		go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", rd, map[string]any{
			"type": types,
			"file": map[string]any{
				"path":    path,
				"content": data,
				"before":  before,
				"after":   after,
				"loguid":  0,
			},
		})
	}

	// var param = interfaces.SystemSettingReq{}
	// if err := api.Context.ShouldBindJSON(&param); err != nil {
	// 	helper.ErrorWith(api.Context, err.Error(), nil)
	// 	return
	// }
	// result, err := system.SystemService.SystemSetting(param)
	// if err != nil {
	// 	helper.ErrorWith(api.Context, err.Error(), nil)
	// 	return
	// }
	helper.Success(api.Context, "result")
}
