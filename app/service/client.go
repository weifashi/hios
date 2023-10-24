package service

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hios/app/model"
	"hios/core"
	"io"
)

var ClientService = clientService{}

type clientService struct{}

// 获取客户端列表
func (ws clientService) List(source string, online string) []model.Client {
	var clients []model.Client
	db := core.DB.Model(&model.Client{})
	if len(source) > 0 && online != "all" {
		db.Where("source = ?", source)
	}
	if len(online) > 0 && online != "all" {
		db.Where("online = ?", online)
	}
	db.Find(&clients)
	return clients
}

// 客户端上线
func (ws clientService) GoLive(uid string, clientIp string, source string) {
	model.ClientModel.UpdateInsert(
		map[string]interface{}{"uid": uid},
		map[string]interface{}{
			"online": 1,
			"source": source,
			"ip":     clientIp,
		},
	)
	//
	clientTableName := core.DBTableName(&model.ClientModel)
	core.DB.Exec(fmt.Sprintf("UPDATE %s SET cc = cc + 1 WHERE ip = ?", clientTableName), clientIp)
}

// 客户端离线
func (ws clientService) Offline(uid string) {
	model.ClientModel.UpdateInsert(
		map[string]interface{}{"uid": uid},
		map[string]interface{}{
			"online": 0,
		},
	)
}

// 创建url
func (ws clientService) CreateUrl(types string, uid string) string {
	url := "/api/v1/ws?type=" + types
	if len(uid) > 0 {
		url = url + "&uid=" + uid
	}
	//
	randomBytes := make([]byte, 32)
	io.ReadFull(rand.Reader, randomBytes)
	hash := md5.Sum(randomBytes)
	sing := hex.EncodeToString(hash[:])
	//
	core.DB.Create(&model.Sing{Sing: sing})
	return url + "&sing=" + sing
}
