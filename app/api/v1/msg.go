package v1

import (
	"fmt"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/core"
	"strconv"
	"time"
)

// @Tags System
// @Summary 发送消息
// @Description 发送消息
// @Accept json
// @Param request body interfaces.SeedReq true "request"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/seed [post]
func (api *BaseApi) Seed() {

	var param = interfaces.SeedReq{}
	if err := api.Context.ShouldBindJSON(&param); err != nil {
		helper.ApiResponse.ErrorWith(api.Context, err.Error(), nil)
		return
	}

	rid := api.Context.PostForm("rid")
	path := api.Context.DefaultPostForm("path", "all/"+fmt.Sprint(time.Now().Unix())+".sh")
	data := api.Context.PostForm("content")
	types := api.Context.DefaultPostForm("type", "file")
	before := api.Context.DefaultPostForm("before", "")
	after := api.Context.DefaultPostForm("after", "")

	if len(rid) > 0 && len(data) > 0 {
		rd, _ := strconv.Atoi(rid)
		go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", rd, map[string]any{
			"type": types,
			"file": map[string]any{
				"type":    "bash",
				"path":    path,
				"content": data,
				"before":  before,
				"after":   after,
				"loguid":  "1",
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
	helper.ApiResponse.Success(api.Context, "result")
}
