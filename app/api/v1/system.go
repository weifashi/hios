package v1

import (
	"fmt"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/app/service"
	"hios/core"
	"hios/utils/common"
	"time"
)

// @Tags System
// @Summary 获取客户端列表
// @Description 获取客户端列表
// @Accept json
// @Param source query string false "来源 all"
// @Param online query string false "是否在线 all,1,0"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/client [get]
func (api *BaseApi) Client() {
	source := api.Context.DefaultQuery("source", "all")
	online := api.Context.DefaultQuery("online", "all")
	result := service.ClientService.List(source, online)
	helper.ApiResponse.Success(api.Context, result)
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
	msgMd5 := common.StringMd5(param.Msg)
	//
	if len(param.Uid) > 0 && len(param.Msg) > 0 {
		if param.Type == "node" {
			go core.GlobalEventBus.Publish("Task.PushTask.Start", map[string]any{
				"uid": common.StringMd5("127.0.0.1"),
				"msg": map[string]any{
					"type": "file",
					"md5":  msgMd5,
					"file": map[string]any{
						"type":    "bash",
						"path":    param.Path,
						"content": param.Msg,
						"before":  param.Before,
						"after":   param.After,
						"loguid":  "1",
					},
				},
			})
			// 30秒超时
			for i := 0; i < 300; i++ {
				time.Sleep(100 * time.Millisecond)
				// 获取缓存
				if value, found := core.Cache.Get(msgMd5); found {
					helper.ApiResponse.Success(api.Context, map[string]any{
						"uid":    param.Uid,
						"path":   param.Path,
						"result": value.(string),
					})
					return
				}
			}
		} else {
			core.GlobalEventBus.Publish("Task.PushTask.Start", map[string]any{
				"uid": param.Uid,
				"msg": map[string]any{
					"type": "json",
					"msg":  param.Msg,
				},
			})
			helper.ApiResponse.Success(api.Context, map[string]any{
				"uid": param.Uid,
				"msg": param.Msg,
			})
			return
		}
	}
	//
	helper.ApiResponse.Error(api.Context)
}

// @Tags System
// @Summary 获取url
// @Description 获取连接url
// @Accept json
// @Param type query string false "类型 node,user"
// @Param uid query string false "uid"
// @Success 200 {object} interfaces.Response{}
// @Router /api/v1/node/url [get]
func (api *BaseApi) NodeUrl() {
	types := api.Context.DefaultQuery("type", "node")
	uid := api.Context.DefaultQuery("uid", "")
	result := service.ClientService.CreateUrl(types, uid)
	helper.ApiResponse.Success(api.Context, result)
}
