package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hios/app/model"
	"hios/core"
	"strings"
	"time"
)

var WebSocketService = webSocketService{}

type webSocketService struct{}

// SaveUser 保存用户
func (ws webSocketService) SaveUser(fd int, userid int) {
	// 一天后过期
	cacheExpiration := 24 * time.Hour
	cacheKeyFD := fmt.Sprintf("User::fd:%d", fd)
	cacheKeyOnline := fmt.Sprintf("User::online:%d", userid)
	core.Cache.Set(cacheKeyFD, "on", cacheExpiration)
	core.Cache.Set(cacheKeyOnline, "on", cacheExpiration)
	// 保存
	key := md5.Sum([]byte(fmt.Sprintf("%d@%d", fd, userid)))
	keyStr := hex.EncodeToString(key[:])
	model.WebSocketModel.UpdateInsert(map[string]interface{}{"key": keyStr}, map[string]interface{}{
		"fd":     fd,
		"userid": userid,
	})
}

// DeleteUser 清除用户
func (ws webSocketService) DeleteUser(fd int) {
	cacheKey := fmt.Sprintf("User::fd:%d", fd)
	core.Cache.Delete(cacheKey)

	var array []string
	db := core.DB.Model(&model.WebSocket{}).Where("fd = ?", fd)
	var webSockets []model.WebSocket
	db.Find(&webSockets)

	for _, webSocket := range webSockets {
		if webSocket.Userid != 0 {
			// 离线时更新会员最后在线时间
			core.DB.Model(&model.User{}).Where("userid = ?", webSocket.Userid).Update("last_at", time.Now().Unix())
			cacheKey := fmt.Sprintf("User::online:%d", webSocket.Userid)
			core.Cache.Delete(cacheKey)
		}
		if strings.HasPrefix(webSocket.Path, "/single/file/") {
			array = append(array, webSocket.Path)
		}
		core.DB.Delete(&webSocket)
	}

	for _, path := range array {
		ws.PushPath(path)
	}
}

// 发送给相同访问状态的会员
func (ws webSocketService) PushPath(path string) {
	// 打印日志
	fmt.Println("推送消息给相同访问状态的会员:", path)
	var userids []int
	core.DB.Model(&model.WebSocket{}).Where("path = ?", path).Pluck("userid", &userids)
	if len(userids) > 0 {
		params := map[string]interface{}{
			"userid": userids,
			"msg": map[string]interface{}{
				"type": "path",
				"data": map[string]interface{}{
					"path":    path,
					"userids": userids,
				},
			},
		}
		// 发送消息任务
		go core.GlobalEventBus.Publish("Task.PushTask.Start", params)
	}
}

// 根据fd获取会员ID
func (ws webSocketService) GetUserid(fd int) int {
	var userid int
	core.DB.Model(&model.WebSocket{}).Where("fd = ?", fd).Pluck("userid", &userid)
	return userid
}
