package service

import (
	"fmt"
	"hios/app/model"
	"hios/core"
)

var ClientService = clientService{}

type clientService struct{}

// 客户端上线
func (ws clientService) GoLive(clientIp string, source string) {
	model.ClientModel.UpdateInsert(
		map[string]interface{}{"ip": clientIp},
		map[string]interface{}{
			"online": 1,
			"source": source,
		},
	)
	//
	clientTableName := core.DBTableName(&model.ClientModel)
	core.DB.Exec(fmt.Sprintf("UPDATE %s SET cc = cc + 1 WHERE ip = ?", clientTableName), clientIp)
}

// 客户端离线
func (ws clientService) Offline(clientIp string) {
	model.ClientModel.UpdateInsert(
		map[string]interface{}{"ip": clientIp},
		map[string]interface{}{
			"online": 0,
		},
	)
}
