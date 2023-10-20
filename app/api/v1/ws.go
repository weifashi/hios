package v1

import (
	"encoding/json"
	"fmt"
	"hios/app/constant"
	"hios/app/helper"
	"hios/app/interfaces"
	"hios/app/service"
	"hios/core"
	"net/http"
	"strconv"
	"sync"
	"time"

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
// @Router /api/v1/ws [get]
func (api *BaseApi) Ws() {
	if api.Context.Request.Header.Get("Upgrade") != "websocket" {
		helper.ApiResponse.ErrorWith(api.Context, constant.ErrNotSupport, nil)
		return
	}
	conn, err := wsUpgrader.Upgrade(api.Context.Writer, api.Context.Request, nil)
	if err != nil {
		helper.ApiResponse.ErrorWith(api.Context, constant.ErrConnFailed, err)
		return
	}

	wsRid++
	client := interfaces.WsClient{
		Conn: conn,
		Type: api.Context.DefaultQuery("type", constant.WsIsUnknown),
		Uid:  api.Context.ClientIP(),
		Rid:  wsRid,
		Ip:   api.Context.ClientIP(),
	}

	//
	if api.Token != "" {
		client.Type = constant.WsIsUser
		if api.Userinfo != nil {
			client.Uid = strconv.Itoa(api.Userinfo.Id)
		}
	}

	// 完成时关闭连接释放资源
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	//
	go func() {
		// 监听连接“完成”事件，其实也可以说丢失事件
		<-api.Context.Done()
		// 客户端离线
		api.wsOfflineClients(client.Rid)
	}()

	// 添加客户端（上线）
	go api.wsOnlineClients(client)

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

		// 缓存执行回调
		if msg.Md5 != "" {
			core.Cache.Set(msg.Md5, msg.Output, 3*time.Second)
			continue
		}
		//
		if msg.Data == nil {
			msg.Data = make(map[string]any)
		}
		if msg.Action == constant.WsHeartbeat {
			// 心跳消息
			go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", int(client.Rid), map[string]any{
				"type": constant.WsHeartbeat,
			})
			continue
		}
		if client.Type == constant.WsIsUser {
			// 用户消息
			api.wsHandleUserMsg(client, msg)
		}
	}

}

// 处理用户消息
func (api *BaseApi) wsHandleUserMsg(client interfaces.WsClient, msg interfaces.WsMsg) {
	fmt.Printf("客户端消息：%v %v\n", client, msg)
	if msg.Action == constant.WsSendMsg {
		// 消息发送
		toType, _ := msg.Data.(map[string]interface{})["type"].(string) // 客户端类型
		toUid, _ := msg.Data.(map[string]interface{})["uid"].(string)   // 发送给谁
		msgData := msg.Data.(map[string]interface{})["data"]            // 消息内容
		msgId := msg.Data.(map[string]interface{})["msgId"]             // 消息ID（用于回调）
		if toUid == "" || msgData == nil {
			return
		}
		// 回调消息
		if msgId != nil {
			go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", int(client.Rid), map[string]interface{}{
				"type":  "receipt",
				"msgId": msgId,
				"data":  map[string]interface{}{},
			})
		}
		sendMsg := interfaces.WsMsg{
			Action: constant.WsSendMsg,
			Data:   msgData,
			Type:   client.Type,
			Uid:    client.Uid,
			Rid:    client.Rid,
		}
		for _, v := range core.WsClients {
			if v.Type == toType && v.Uid == toUid {
				go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", toUid, sendMsg)
			}
		}
	}
}

// 客户端上线
func (api *BaseApi) wsOnlineClients(client interfaces.WsClient) {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	//
	for _, v := range core.WsClients {
		if v.Rid == client.Rid {
			return
		}
	}
	core.WsClients = append(core.WsClients, client)

	// 保存用户
	service.WebSocketService.SaveUser(int(client.Rid), client.Uid)

	// 客户端上线
	service.ClientService.GoLive(client.Ip, client.Type)

	// 发送open事件
	go core.GlobalEventBus.Publish("Task.PushTask.PushMsg", int(client.Rid), map[string]interface{}{
		"type": "open",
		"data": map[string]interface{}{
			"fd": int(client.Rid),
		},
	})

	// 通知上线
	go core.GlobalEventBus.Publish("Task.LineTask.Start", client.Uid, true)

	// 推送离线时收到的消息
	go core.GlobalEventBus.Publish("Task.PushTask.Start", "RETRY::"+client.Uid)
}

// 客户端离线
func (api *BaseApi) wsOfflineClients(rid int32) {
	wsMutex.Lock()
	defer wsMutex.Unlock()
	for k, client := range core.WsClients {
		if client.Rid == rid {
			core.WsClients = append(core.WsClients[:k], core.WsClients[k+1:]...)
			_ = client.Conn.Close()
			// 通知离线
			go core.GlobalEventBus.Publish("Task.LineTask.Start", client.Uid, false)
			// 清除用户
			service.WebSocketService.DeleteUser(int(client.Rid))
			// 客户端离线
			service.ClientService.Offline(client.Ip)
			break
		}
	}
}
