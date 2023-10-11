package wsc

import (
	"encoding/json"
	"fmt"
	"hios/config"
	"hios/utils/cmd"
	"hios/utils/common"
	"hios/utils/logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/togettoyou/wsc"
)

var (
	ws      *wsc.Wsc
	FileMd5 sync.Map
	workDir = "work"
)

type msgModel struct {
	Type    string    `json:"type"`
	Content string    `json:"content"`
	File    fileModel `json:"file"`
	Cmd     cmdModel  `json:"cmd"`
}

type fileModel struct {
	Type    string `json:"type"`
	Path    string `json:"path"`
	Before  string `json:"before"`
	After   string `json:"after"`
	Content string `json:"content"`
	Loguid  string `json:"loguid"`
}

type cmdModel struct {
	Log      bool   `json:"log"`
	Callback string `json:"callback"`
	Content  string `json:"content"`
	Loguid   string `json:"loguid"`
}

type sendModel struct {
	Type   string      `json:"type"`
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type callModel struct {
	Callback string `json:"callback"`
	Output   string `json:"output"`
	Err      string `json:"err"`
}

// WorkStart Work开始
func WorkStart() {
	//
	initWorkDir()
	//
	logger.SetLogger(`{"File":{"filename":"work/wsc.log","level":"TRAC","daily":true,"maxlines":100000,"maxsize":10,"maxdays":3,"append":true,"permit":"0660"}}`)
	//
	origin := strings.Replace(config.CONF.System.WssUrl, "https://", "wss://", 1)
	origin = strings.Replace(origin, "http://", "ws://", 1)
	//
	done := make(chan bool)
	//
	ws = wsc.New(origin)
	// 自定义配置
	ws.SetConfig(&wsc.Config{
		WriteWait:         2 * time.Second,
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
		logger.Error("[ws] connect error: ", err.Error())
	})
	ws.OnDisconnected(func(err error) {
		logger.Error("[ws] disconnected: ", err.Error())
	})
	ws.OnClose(func(code int, text string) {
		logger.Info("[ws] close: ", code, text)
		done <- true
	})
	ws.OnTextMessageSent(func(message string) {
		if !strings.HasPrefix(message, "r:") {
			logger.Debug("[ws] text message sent: ", message)
		}
	})
	ws.OnBinaryMessageSent(func(data []byte) {
		logger.Info("[ws] binary message sent: ", string(data))
	})
	ws.OnSentError(func(err error) {
		logger.Info("[ws] sent error: ", err.Error())
	})
	ws.OnPingReceived(func(appData string) {
		logger.Info("[ws] ping received: ", appData)
	})
	ws.OnPongReceived(func(appData string) {
		logger.Info("[ws] pong received: ", appData)
	})
	ws.OnTextMessageReceived(func(message string) {
		logger.Info("[ws] text message received: ", message)
		// if strings.HasPrefix(message, "r:") {
		// 	message = xrsa.Decrypt(message[2:], nodePublic, nodePrivate) // 判断数据解密
		// }
		handleMessageReceived(message)
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

// 初始化工作目录
func initWorkDir() {
	err := common.Mkdir(workDir, 0777)
	if err != nil {
		logger.Error(fmt.Sprintf("[start] failed to create log dir: %s\n", err.Error()))
		os.Exit(1)
	}
}

// 处理消息
func handleMessageReceived(message string) {
	var data msgModel
	if ok := json.Unmarshal([]byte(message), &data); ok == nil {
		if data.Type == "file" {
			// 保存文件
			handleMessageFile(data.File, false)
		} else if data.Type == "cmd" {
			// 执行命令
			output, err := handleMessageCmd(data.Cmd.Content, data.Cmd.Log, data.Cmd.Loguid)
			if len(data.Cmd.Callback) > 0 {
				cmderr := ""
				if err != nil {
					cmderr = err.Error()
				}
				callData := &callModel{
					Callback: data.Cmd.Callback,
					Output:   output,
					Err:      cmderr}
				sendMessage := formatSendMsg("cmd", callData)
				err = ws.SendTextMessage(sendMessage)
				if err != nil {
					logger.Debug("[cmd] send callback error: %s", err)
				}
			}
		}
	}
}

// 保存文件或运行文件
func handleMessageFile(fileData fileModel, force bool) {
	var err error
	var output string
	if !strings.HasPrefix(fileData.Path, "/") {
		fileData.Path = fmt.Sprintf("%s/%s", workDir, fileData.Path)
	}
	fileDir := filepath.Dir(fileData.Path)
	if !common.Exists(fileDir) {
		err = os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			logger.Error("#%s# [file] mkdir error: '%s' %s", fileData.Loguid, fileDir, err)
			return
		}
	}
	fileContent := fileData.Content
	if fileContent == "" {
		logger.Warn("#%s# [file] empty: %s", fileData.Loguid, fileData.Path)
		return
	}
	//
	fileKey := common.StringMd5(fileData.Path)
	contentKey := common.StringMd5(fileContent)
	if !force {
		md5Value, _ := FileMd5.Load(fileKey)
		if md5Value != nil && md5Value.(string) == contentKey {
			logger.Debug("#%s# [file] same: %s", fileData.Loguid, fileData.Path)
			return
		}
	}
	FileMd5.Store(fileKey, contentKey)
	//
	if len(fileData.Before) > 0 {
		beforeFile := fmt.Sprintf("%s.before", fileData.Path)
		err = os.WriteFile(beforeFile, []byte(fileData.Before), 0666)
		if err != nil {
			logger.Error("#%s# [before] write before error: '%s' %s", fileData.Loguid, beforeFile, err)
			return
		}
		logger.Info("#%s# [before] start: '%s'", fileData.Loguid, beforeFile)
		_, _ = cmd.Bash("-c", fmt.Sprintf("chmod +x %s", beforeFile))
		output, err = cmd.Bash(beforeFile)
		if err != nil {
			logger.Error("#%s# [before] error: '%s' %s %s", fileData.Loguid, beforeFile, err, output)
		} else {
			logger.Info("#%s# [before] success: '%s'", fileData.Loguid, beforeFile)
		}
	}
	//
	err = os.WriteFile(fileData.Path, []byte(fileContent), 0666)
	if err != nil {
		logger.Error("#%s# [file] write error: '%s' %s", fileData.Loguid, fileData.Path, err)
		return
	}
	if common.InArray(fileData.Type, []string{"bash", "cmd", "exec"}) {
		logger.Info("#%s# [bash] start: '%s'", fileData.Loguid, fileData.Path)
		_, _ = cmd.Bash("-c", fmt.Sprintf("chmod +x %s", fileData.Path))
		output, err = cmd.Bash(fileData.Path)
		if err != nil {
			logger.Error("#%s# [bash] error: '%s' %s %s", fileData.Path, err, output)
		} else {
			logger.Info("#%s# [bash] success: '%s'", fileData.Loguid, fileData.Path)
		}
	} else if fileData.Type == "sh" {
		logger.Info("#%s# [sh] start: '%s'", fileData.Loguid, fileData.Path)
		_, _ = cmd.Cmd("-c", fmt.Sprintf("chmod +x %s", fileData.Path))
		output, err = cmd.Cmd(fileData.Path)
		if err != nil {
			logger.Error("#%s# [sh] error: '%s' %s %s", fileData.Path, err, output)
		} else {
			logger.Info("#%s# [sh] success: '%s'", fileData.Loguid, fileData.Path)
		}
	} else if fileData.Type == "yml" {
		logger.Info("#%s# [yml] start: '%s'", fileData.Loguid, fileData.Path)
		output, err = cmd.Cmd("-c", fmt.Sprintf("cd %s && docker-compose up -d --remove-orphans", fileDir))
		if err != nil {
			logger.Error("#%s# [yml] error: '%s' %s %s", fileData.Loguid, fileData.Path, err, output)
		} else {
			logger.Info("#%s# [yml] success: '%s'", fileData.Loguid, fileData.Path)
		}
	} else if fileData.Type == "nginx" {
		logger.Info("#%s# [nginx] start: '%s'", fileData.Loguid, fileData.Path)
		output, err = cmd.Cmd("-c", "nginx -s reload")
		if err != nil {
			logger.Error("#%s# [nginx] error: '%s' %s %s", fileData.Loguid, fileData.Path, err, output)
		} else {
			logger.Info("#%s# [nginx] success: '%s'", fileData.Loguid, fileData.Path)
		}
	}
	//
	if len(fileData.After) > 0 {
		afterFile := fmt.Sprintf("%s.after", fileData.Path)
		err = os.WriteFile(afterFile, []byte(fileData.After), 0666)
		if err != nil {
			logger.Error("#%s# [after] write after error: '%s' %s", fileData.Loguid, afterFile, err)
			return
		}
		logger.Info("#%s# [after] start: '%s'", fileData.Loguid, afterFile)
		_, _ = cmd.Bash("-c", fmt.Sprintf("chmod +x %s", afterFile))
		output, err = cmd.Bash(afterFile)
		if err != nil {
			logger.Error("#%s# [after] error: '%s' %s %s", fileData.Loguid, afterFile, err, output)
		} else {
			logger.Info("#%s# [after] success: '%s'", fileData.Loguid, afterFile)
		}
	}
}

// 运行自定义脚本
func handleMessageCmd(cmds string, addLog bool, loguid string) (string, error) {
	output, err := cmd.Cmd("-c", cmds)
	if addLog {
		if err != nil {
			logger.Error("#%s# [cmd] error: '%s' %s; output: '%s'", loguid, cmds, err, output)
		} else {
			logger.Info("#%s# [cmd] success: '%s'", loguid, cmds)
		}
	}
	return output, err
}

// 格式化要发送的消息
func formatSendMsg(action string, data interface{}) string {
	sendData := &sendModel{Type: "node", Action: action, Data: data}
	sendRes, sendErr := json.Marshal(sendData)
	if sendErr != nil {
		return ""
	}
	msg := string(sendRes)
	// if len(serverPublic) > 0 {
	// 	return fmt.Sprintf("r:%s", xrsa.Encrypt(msg, serverPublic))
	// } else {
	return msg
	// }
}
