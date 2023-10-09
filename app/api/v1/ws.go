package v1

import (
	"encoding/json"
	"fmt"
	"hios/app/constant"
	"hios/app/helper"
	"hios/app/interfaces"
	service "hios/app/service/websocket"
	"hios/core"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	wsRid      int32 = 0
	wsMutex          = sync.Mutex{}
	wsUpgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// @Tags Websocket
// @Summary Websocket ws
// @Description 请使用ws连接
// @Accept json
// @Param request query interfaces.WebsocketReq true "request"
// @Router /ws [get]
func (api *BaseApi) Ws() {
	if api.Context.Request.Header.Get("Upgrade") != "websocket" {
		helper.ErrorWith(api.Context, constant.ErrNotSupport, nil)
		return
	}
	conn, err := wsUpgrader.Upgrade(api.Context.Writer, api.Context.Request, nil)
	if err != nil {
		helper.ErrorWith(api.Context, constant.ErrConnFailed, err)
		return
	}
	wsRid++
	client := interfaces.WsClient{
		Conn: conn,
		Type: interfaces.WsIsUnknown,
		Uid:  0,
		Rid:  wsRid,
	}
	if api.Token != "" {
		client.Type = interfaces.WsIsUser
		client.Uid = int32(api.Userinfo.Userid)
	}
	// 完成时关闭连接释放资源
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-api.Context.Done()
		// 客户端离线
		api.wsOfflineClients(client.Rid)
	}()
	// 添加客户端（上线）
	api.wsOnlineClients(client)
	sendMsg, _ := json.Marshal(interfaces.WsMsg{
		Action: interfaces.WsOnline,
		Data:   map[string]any{"own": 1},
		Type:   client.Type,
		Uid:    client.Uid,
		Rid:    client.Rid,
	})
	_ = conn.WriteMessage(websocket.TextMessage, sendMsg)
	// 循环读取客户端发送的消息
	for {
		// 读取客户端发送过来的消息，如果没发就会一直阻塞住
		_, message, err := conn.ReadMessage()
		if err != nil {
			api.wsOfflineClients(client.Rid)
			break
		}
		var msg interfaces.WsMsg
		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		}
		if msg.Data == nil {
			msg.Data = make(map[string]any)
		}
		if msg.Action == interfaces.WsHeartbeat {
			// 心跳消息
			sendMsg, _ = json.Marshal(map[string]any{
				"type": interfaces.WsHeartbeat,
			})
			_ = conn.WriteMessage(websocket.TextMessage, sendMsg)
			continue
		}
		if client.Type == interfaces.WsIsUser {
			// 用户消息
			api.wsHandleUserMsg(client, msg)
		}
	}
}

// 处理用户消息
func (api *BaseApi) wsHandleUserMsg(client interfaces.WsClient, msg interfaces.WsMsg) {
	fmt.Printf("客户端消息：%v %v\n", client, msg)
	var replyMsg []byte
	if msg.Action == interfaces.WsSendMsg {
		// 消息发送
		toType, _ := msg.Data.(map[string]interface{})["type"].(string) // 客户端类型
		toUid, _ := msg.Data.(map[string]interface{})["uid"].(float64)  // 发送给谁
		msgData := msg.Data.(map[string]interface{})["data"]            // 消息内容
		msgId := msg.Data.(map[string]interface{})["msgId"]             // 消息ID（用于回调）
		if toUid == 0 || msgData == nil {
			return
		}
		// 回调消息
		if msgId != nil {
			core.GlobalEventBus.Publish("Task.PushTask.PushMsg", int(client.Rid), map[string]interface{}{
				"type":  "receipt",
				"msgId": msgId,
				"data":  map[string]interface{}{},
			})
		}
		sendMsg, _ := json.Marshal(interfaces.WsMsg{
			Action: interfaces.WsSendMsg,
			Data:   msgData,
			Type:   client.Type,
			Uid:    client.Uid,
			Rid:    client.Rid,
		})
		for _, v := range interfaces.WsClients {
			if v.Type == toType && v.Uid == int32(toUid) {
				_ = v.Conn.WriteMessage(websocket.TextMessage, sendMsg)
			}
		}
	}
	if replyMsg != nil {
		_ = client.Conn.WriteMessage(websocket.TextMessage, replyMsg)
	}
}

// 客户端上线
func (api *BaseApi) wsOnlineClients(client interfaces.WsClient) {
	for _, v := range interfaces.WsClients {
		if v.Rid == client.Rid {
			return
		}
	}
	wsMutex.Lock()
	interfaces.WsClients = append(interfaces.WsClients, client)
	wsMutex.Unlock()
	// 处理上线事件
	fmt.Printf("处理上线事件：%v \n", client)
	// 保存用户
	service.WebSocketService.SaveUser(int(client.Rid), int(client.Uid))
	// 发送open事件
	core.GlobalEventBus.Publish("Task.PushTask.PushMsg", int(client.Rid), map[string]interface{}{
		"type": "open",
		"data": map[string]interface{}{
			"fd": int(client.Rid),
		},
	})

	// 通知上线
	core.GlobalEventBus.Publish("Task.LineTask.Start", int(client.Uid), true)

	// 推送离线时收到的消息
	core.GlobalEventBus.Publish("Task.PushTask.Start", "RETRY::"+strconv.Itoa(int(client.Uid)))
}

// 客户端离线
func (api *BaseApi) wsOfflineClients(rid int32) {
	for k, client := range interfaces.WsClients {
		if client.Rid == rid {
			wsMutex.Lock()
			interfaces.WsClients = append(interfaces.WsClients[:k], interfaces.WsClients[k+1:]...)
			_ = client.Conn.Close()
			wsMutex.Unlock()
			// 处理离线事件
			fmt.Printf("处理离线事件：%v \n", client)
			// 通知离线
			core.GlobalEventBus.Publish("Task.LineTask.Start", int(client.Uid), false)
			// 清除用户
			service.WebSocketService.DeleteUser(int(client.Rid))
			break
		}
	}
}
