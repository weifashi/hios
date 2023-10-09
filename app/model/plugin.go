package model

import "hios/core"

type Plugin struct {
	core.BaseIdModels
	Key  string `gorm:"type:varchar(50);comment:插件key" json:"key"`
	Ver  string `gorm:"type:varchar(25);comment:版本号" json:"ver"`
	Data string `gorm:"comment:安装数据" json:"data"`
	core.BaseAtModels
}

var PluginModel = Plugin{}

// 根据code获取信息
func (m Plugin) GetPluginByCode(id string) (*Plugin, error) {
	err := core.DB.Where("id=?", id).Find(&m).Error
	return &m, err
}
