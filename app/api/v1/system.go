package v1

import (
	"fmt"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/core"
	"hios/utils/common"
	"time"
)

// @Tags System
// @Summary 获取客户端列表
// @Description 获取客户端列表
// @Accept json
// @Param request body interfaces.SeedReq true "request"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/client [get]
func (api *BaseApi) Client() {
	var param = interfaces.SeedReq{
		Path: "all/" + fmt.Sprint(time.Now().Unix()) + ".sh",
	}
	helper.ApiRequest.ShouldBindAll(api.Context, &param)
	//
	md5 := common.StringMd5(param.Cmd)
	//
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		// 获取缓存
		if value, found := core.Cache.Get(md5); found {
			helper.ApiResponse.Success(api.Context, map[string]any{
				"rid":    param.Rid,
				"path":   param.Path,
				"result": value.(string),
			})
			return
		}
	}
	//
	helper.ApiResponse.Error(api.Context)
}

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

		go core.GlobalEventBus.Publish("Task.PushTask.Start", map[string]any{
			"uid": "127.0.0.1",
			"msg": map[string]any{
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
			},
		})

		// rd, _ := strconv.Atoi(param.Rid)
		// go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", rd, map[string]any{
		// 	"type": "file",
		// 	"md5":  md5,
		// 	"file": map[string]any{
		// 		"type":    "bash",
		// 		"path":    param.Path,
		// 		"content": param.Cmd,
		// 		"before":  param.Before,
		// 		"after":   param.After,
		// 		"loguid":  "1",
		// 	},
		// })
	}
	// 30秒超时
	for i := 0; i < 300; i++ {
		time.Sleep(100 * time.Millisecond)
		// 获取缓存
		if value, found := core.Cache.Get(md5); found {
			helper.ApiResponse.Success(api.Context, map[string]any{
				"rid":    param.Rid,
				"path":   param.Path,
				"result": value.(string),
			})
			return
		}
	}
	//
	helper.ApiResponse.Error(api.Context)
}
