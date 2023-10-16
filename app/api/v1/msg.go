package v1

import (
	"fmt"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/core"
	"hios/utils/common"
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
	var param = interfaces.SeedReq{
		Path: "all/" + fmt.Sprint(time.Now().Unix()) + ".sh",
	}
	helper.ApiRequest.ShouldBindAll(api.Context, &param)
	//
	md5 := common.StringMd5(param.Cmd)
	//
	if len(param.Rid) > 0 && len(param.Cmd) > 0 {
		rd, _ := strconv.Atoi(param.Rid)
		go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", rd, map[string]any{
			"type": "file",
			"md5":  md5,
			"file": map[string]any{
				"type":    "bash",
				"path":    param.Path,
				"content": param.Cmd,
				"before":  param.Before,
				"after":   param.After,
				"loguid":  "1",
			},
		})
	}
	// 30秒超时
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		// 获取缓存
		if value, found := core.Cache.Get(md5); found {
			param.Result = value.(string)
			helper.ApiResponse.Success(api.Context, param)
			return
		}
	}
	//
	helper.ApiResponse.Error(api.Context)
}
