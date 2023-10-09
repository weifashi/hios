package task

import (
	"context"
	"hios/app/model"
	"hios/core"
	"log"
	"strconv"
	"time"
)

// 通知用户上线/离线状态的任务。
var LineTask = lineTask{}

type lineTask struct{}

// 初始添加事件总线的订阅
func init() {
	core.GlobalEventBus.SubscribeAsync("Task.LineTask.Start", new(lineTask).Start, true)
}

// userid: 用户的 ID。
// online: 用户的上线/离线状态 true: 上线，false: 离线。
func (t lineTask) Start(userid int, online bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var endPush []map[string]interface{}
	var websockets []model.WebSocket
	if err := core.DB.WithContext(ctx).Where("userid != ?", userid).Find(&websockets).Error; err != nil {
		log.Printf("查询 WebSocket 失败: %v", err)
		return
	}
	for _, ws := range websockets {
		fd, err := strconv.Atoi(ws.Fd)
		if err != nil {
			log.Printf("将 fd 转换为 int 失败: %v", err)
			continue
		}
		endPush = append(endPush, map[string]interface{}{
			"fd": fd,
			"msg": map[string]interface{}{
				"type": "line",
				"data": map[string]interface{}{
					"userid": userid,
					"online": online,
				},
			},
		})
	}
	if len(endPush) == 0 {
		return
	}
	// 推送信息
	PushTask.Start(endPush)
}
