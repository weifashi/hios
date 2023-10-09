package wsc

import (
	"hios/config"
	"hios/utils/logger"
	"strings"
	"time"

	"github.com/togettoyou/wsc"
)

var (
	ws *wsc.Wsc
)

// WorkStart Work开始
func WorkStart() {
	//
	logger.SetLogger(`{"File":{"filename":"wsc.log","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`)
	//
	origin := strings.Replace(config.CONF.System.WssUrl, "https://", "wss://", 1)
	origin = strings.Replace(origin, "http://", "ws://", 1)
	//
	done := make(chan bool)
	//
	ws = wsc.New(origin)
	// 自定义配置
	ws.SetConfig(&wsc.Config{
		WriteWait:         10 * time.Second,
		MaxMessageSize:    512 * 1024, // 512KB
		MinRecTime:        3 * time.Second,
		MaxRecTime:        30 * time.Second,
		RecFactor:         10,
		MessageBufferSize: 1024,
	})
	// 设置回调处理
	ws.OnConnected(func() {
		logger.Debug("[ws] connected: ", ws.WebSocket.Url)
		logger.SetWebsocket(ws)
		// onConnected()
	})
	ws.OnConnectError(func(err error) {
		logger.Debug("[ws] connect error: ", err.Error())
		// switchServer()
	})
	ws.OnDisconnected(func(err error) {
		// logger.Debug("[ws] disconnected: ", err.Error())
		// switchServer()
	})
	ws.OnClose(func(code int, text string) {
		logger.Debug("[ws] close: ", code, text)
		done <- true
	})
	ws.OnTextMessageSent(func(message string) {
		if !strings.HasPrefix(message, "r:") {
			logger.Debug("[ws] text message sent: ", message)
		}
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		logger.Debug("[ws] binary message sent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		logger.Debug("[ws] sent error: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		logger.Debug("[ws] ping received: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		logger.Debug("[ws] pong received: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		// if strings.HasPrefix(message, "r:") {
		// 	message = xrsa.Decrypt(message[2:], nodePublic, nodePrivate) // 判断数据解密
		// } else {
		logger.Info("[ws] text message received: ", message)
		// }
		// handleMessageReceived(message)
	})
	ws.OnBinaryMessageReceived(func(data []byte) {
		logger.Debug("[ws] binary message received: ", string(data))
	})
	// 开始连接
	go ws.Connect()
	//
	for {
		select {
		case <-done:
			return
		}
	}
}
