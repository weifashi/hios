package task

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hios/app/interfaces"
	"hios/app/model"
	"hios/core"
	"hios/utils/common"
	"log"
	"reflect"
	"strconv"
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
func (t pushTask) Start(param interface{}, retryOffline ...bool) {
	//
	retryOfflines := true
	if len(retryOffline) > 0 {
		retryOfflines = retryOffline[0]
	}
	if p, ok := param.(string); ok {
		if strings.HasPrefix(p, "PUSH::") {
			// 推送缓存
			if v, _ := core.Cache.Get(p); v != nil {
				if m, ok := v.(map[string]interface{}); ok && m["fd"] != nil {
					t.push([]map[string]interface{}{m}, retryOfflines, "", 0)
				}
			}
		} else if strings.HasPrefix(p, "RETRY::") {
			// 根据会员ID推送离线时收到的消息
			userid := strings.TrimPrefix(p, "RETRY::")
			p := t.sendTmpMsgForUserid(common.StringToInt(userid))
			t.push(p, retryOfflines, "", 0)
		}
	} else if p, ok := param.(map[string]interface{}); ok {
		t.push([]map[string]interface{}{p}, retryOfflines, "", 0)
	} else if p, ok := param.([]map[string]interface{}); ok {
		t.push(p, retryOfflines, "", 0)
	}
}

// 根据会员ID推送离线时收到的消息。
func (t pushTask) sendTmpMsgForUserid(userid int) []map[string]interface{} {
	if userid == 0 {
		return nil
	}

	var endPush []map[string]interface{}
	var tmpMsgs []model.WebSocketTmpMsg
	if err := core.DB.Where("create_id = ? AND send = 0 AND created_at > ?", userid, time.Now().Add(-time.Minute).Unix()).
		Order("id").Find(&tmpMsgs).Error; err != nil {
		log.Printf("查询 WebSocketTmpMsg 失败: %v", err)
		return nil
	}
	fmt.Println("查询到的离线消息：" + strconv.Itoa(len(tmpMsgs)))
	for _, item := range tmpMsgs {
		msg, err := common.StrToMap(item.Msg)
		if err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}
		endPush = append(endPush, map[string]interface{}{
			"tmpMsgId": item.Id,
			"userid":   item.CreateId,
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
func (t pushTask) addTmpMsg(userFail []int, msg map[string]interface{}) {
	for _, uid := range userFail {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Printf("序列化消息失败: %v\n", err)
			continue
		}
		inArray := model.WebSocketTmpMsg{
			Md5:      fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d-%s", uid, msgBytes))))[8:24],
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
	if err := core.DB.Model(&model.WebSocketTmpMsg{}).Where("md5 = ?", md5).Count(&count).Error; err != nil {
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
// retryOffline: 如果会员不在线，等上线后继续发送
// key: 延迟推送 key 依据，留空立即推送（延迟推送时发给同一人同一种消息类型只发送最新的一条）
// delay: 延迟推送时间，默认：1 秒（key 填写时有效）
func (t pushTask) push(lists []map[string]interface{}, retryOffline bool, key string, delay int) {
	//
	if len(lists) == 0 {
		return
	}
	for _, item := range lists {
		if !common.IsKind(item, reflect.Map) || len(item) == 0 {
			continue
		}
		userid := item["userid"]
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
		offlineUser := []int{}
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
		if userid != nil {
			if !common.IsKind(userid, reflect.Map) {
				userid = []int{userid.(int)}
			}

			for _, uid := range userid.([]int) {
				var row []int
				core.DB.Model(&model.WebSocket{}).Where("userid = ?", userid).Pluck("fd", &row)
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
			if key != "" {
				key = "PUSH::" + strconv.Itoa(fid) + ":" + msgType + ":" + key
				core.Cache.Set(key, map[string]interface{}{
					"fd":       fid,
					"ignoreFd": ignoreFd,
					"msg":      msg,
				}, 24*time.Hour)
				// todo 延迟推送
			} else {
				// 发送消息
				t.PushMsg(fid, msg)
				if tmpMsgId > 0 {
					fmt.Println("更新离线消息状态：" + strconv.Itoa(tmpMsgId))
					core.DB.Model(&model.WebSocketTmpMsg{}).Where("id = ?", tmpMsgId).Update("send", 1)
				}
			}
		}
		// 记录不在线的
		if retryOffline && tmpMsgId == 0 {
			offlineUser = common.UniqueInt(offlineUser)
			t.addTmpMsg(offlineUser, msg)
		}
	}
}

// 推送消息
func (t pushTask) PushMsg(fd int, msg interface{}) {
	//
	for _, v := range core.WsClients {
		if v.Rid == int32(fd) {
			log.Println("推送消息给fd:", v.Rid)
			log.Printf("推送消息内容:%s", common.StructToJson(msg))
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
	// for data := range sendChan {
	if v.Conn != nil && v.Conn.UnderlyingConn() != nil {
		if err := v.Conn.WriteMessage(websocket.TextMessage, msgJSON); err != nil {
			if !websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err.Error())
			}
		}
	}
	// }

}
