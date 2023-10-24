package service

import (
	"fmt"
	"hios/app/model"
	"hios/core"
)

var ClientService = clientService{}

type clientService struct{}

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
