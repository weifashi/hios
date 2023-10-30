package task

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hios/app/interfaces"
	"hios/app/model"
	"hios/core"
	"hios/utils/common"
	"hios/utils/logger"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var pushTaskMutex sync.Mutex

// 推送任务
var PushTask = pushTask{}

type pushTask struct{}

// 初始添加事件总线的订阅
func init() {
	core.GlobalEventBus.SubscribeAsync("Task.PushTask.Start", new(pushTask).Start, true)
	core.GlobalEventBus.SubscribeAsync("Task.PushTask.PushMsg", new(pushTask).PushMsg, true)
}

// 开始
// param 推送的消息
// delay 延时（秒）
func (t pushTask) Start(param interface{}, delay ...int64) {
	//
	var d int64
	if len(delay) > 0 {
		d = delay[0]
	}
	//
	if p, ok := param.(string); ok {
		if strings.HasPrefix(p, "RETRY::") {
			p := t.sendTmpMsgForUserid(strings.TrimPrefix(p, "RETRY::"))
			t.push(p, d)
		}
	} else if p, ok := param.(map[string]interface{}); ok {
		t.push([]map[string]interface{}{p}, d)
	} else if p, ok := param.([]map[string]interface{}); ok {
		t.push(p, d)
	}
}

// 根据会员ID推送离线时收到的消息。
func (t pushTask) sendTmpMsgForUserid(uid string) []map[string]interface{} {
	if uid == "" {
		return nil
	}

	var endPush []map[string]interface{}
	var tmpMsgs []model.WebSocketTmpMsg
	if err := core.DB.Where("create_id = ? AND send = 0 AND created_at > ?", uid, time.Now().Add(-time.Minute).Unix()).
		Order("id").Find(&tmpMsgs).Error; err != nil {
		log.Printf("查询 WebSocketTmpMsg 失败: %v", err)
		return nil
	}
	for _, item := range tmpMsgs {
		msg, err := common.StrToMap(item.Msg)
		if err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}
		endPush = append(endPush, map[string]interface{}{
			"tmpMsgId": item.Id,
			"uid":      item.CreateId,
			"msg":      msg,
		})
	}
	if len(endPush) == 0 {
		return nil
	}
	return endPush
}

// userFail: 发送失败的用户 ID 列表。
// msg: 发送的消息。
func (t pushTask) addTmpMsg(userFail []string, msg map[string]interface{}) {
	for _, uid := range userFail {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("序列化消息失败: %v\n", err)
			continue
		}
		inArray := model.WebSocketTmpMsg{
			Md5:      fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%s", uid, msgBytes))))[8:24],
			Msg:      string(msgBytes),
			Send:     0,
			CreateId: uid,
		}
		if !t.checkMd5(inArray.Md5) {
			if err := t.insertWebSocketTmpMsg(inArray); err != nil {
				fmt.Printf("插入 WebSocketTmpMsg 失败: %v\n", err)
			}
		}
	}
}

// 检查指定的 MD5 值是否已经存在。
func (t pushTask) checkMd5(md5 string) bool {
	var count int64
	if err := core.DB.Model(&model.WebSocketTmpMsg{}).Where("md5 = ? and send = 0", md5).Count(&count).Error; err != nil {
		fmt.Printf("查询 WebSocketTmpMsg 失败: %v\n", err)
		return false
	}
	return count > 0
}

// 插入 WebSocketTmpMsg 记录。
func (t pushTask) insertWebSocketTmpMsg(inArray model.WebSocketTmpMsg) error {
	return core.DB.Create(&inArray).Error
}

// lists: 消息列表
// key: 延迟推送 key 依据，留空立即推送（延迟推送时发给同一人同一种消息类型只发送最新的一条）
// delay: 延迟推送时间，默认：1 秒（key 填写时有效）
func (t pushTask) push(lists []map[string]interface{}, delay ...int64) {
	//
	if len(lists) == 0 {
		return
	}
	for _, item := range lists {
		if !common.IsKind(item, reflect.Map) || len(item) == 0 {
			continue
		}
		uid := item["uid"]
		fd := item["fd"]
		ignoreFd := item["ignoreFd"]
		msg, _ := item["msg"].(map[string]interface{})
		tmpMsgId, _ := item["tmpMsgId"].(int)
		if !common.IsKind(msg, reflect.Map) || len(msg) == 0 {
			continue
		}
		msgType := msg["type"].(string)
		if msgType == "" {
			continue
		}
		// 发送对象
		offlineUser := []string{}
		array := []int{}
		//
		if fd != nil {
			if common.IsKind(fd, reflect.Map) {
				for _, f := range fd.([]interface{}) {
					array = append(array, f.(int))
				}
			} else {
				array = append(array, fd.(int))
			}
		}
		//
		if uid != nil {
			if !common.IsKind(uid, reflect.Map) {
				uid = []string{uid.(string)}
			}
			for _, uid := range uid.([]string) {
				var row []int
				core.DB.Model(&model.WebSocket{}).Order("id DESC").Where("uid = ?", uid).Pluck("fd", &row)
				if len(row) > 0 {
					array = append(array, row...)
				} else {
					offlineUser = append(offlineUser, uid)
				}
			}
		}
		//
		if ignoreFd != nil {
			if !common.IsKind(ignoreFd, reflect.Map) {
				ignoreFd = []interface{}{ignoreFd}
			}
		}
		// 开始发送
		for _, fid := range array {
			if ignoreFd != nil {
				if common.InSlice(fid, ignoreFd.([]interface{})) {
					continue
				}
			}
			if len(delay) != 0 && delay[0] != 0 {
				go func(fid int, msg map[string]interface{}, delay int64) {
					time.Sleep(time.Duration(delay) * time.Second)
					t.PushMsg(fid, msg)
				}(fid, msg, delay[0])
			} else {
				// 发送消息
				t.PushMsg(fid, msg)
				if tmpMsgId > 0 {
					core.DB.Where("id = ?", tmpMsgId).Delete(&model.WebSocketTmpMsg{})
				}
			}
		}
		// 记录不在线的
		if tmpMsgId == 0 {
			t.addTmpMsg(common.UniqueStr(offlineUser), msg)
		}
	}
}

// 推送消息
func (t pushTask) PushMsg(fd int, msg interface{}) {
	for _, v := range core.WsClients {
		if v.Rid == int32(fd) {
			msgJSON, err := json.Marshal(msg)
			if err != nil {
				log.Println("Failed to convert message to JSON:", err)
				continue
			}
			go t.pushWriteMessage(v, msgJSON)
		}
	}
}

// 写入消息
func (t pushTask) pushWriteMessage(v interfaces.WsClient, msgJSON []byte) {
	pushTaskMutex.Lock()
	defer pushTaskMutex.Unlock()
	if v.Conn != nil && v.Conn.UnderlyingConn() != nil {
		//
		logger.SetLogger(`{"File":{"filename":"./work/logs/wss.log","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`)
		logger.Info("[%s] %s", "wss-send", msgJSON)
		//
		if err := v.Conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err.Error())
			}
		}
	}
}
