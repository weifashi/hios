package model

import (
	"encoding/json"
	"hios/core"
	"sync"
)

// Setting 设置实例
type Setting struct {
	core.BaseIdModels
	Name    string `gorm:"type:varchar(50);comment:参数名称" json:"name"`
	Desc    string `gorm:"type:varchar(50);comment:参数描述、备注" json:"desc"`
	Setting string `gorm:"comment:设置的json数据" json:"setting"`
	core.BaseAtModels
}

var (
	SettingSystemKey  = "system"         //系统设置key
	SettingEmailKey   = "emailSetting"   //邮箱设置key
	SettingAppPushKey = "appPushSetting" //App推送设置key

	cache   = make(map[string][]byte) // 缓存
	cacheMu sync.RWMutex              // 读写锁
)

var SettingModel = Setting{}

// GetSetting 获取设置
func (m Setting) GetSetting(name string) (map[string]interface{}, error) {
	cacheMu.RLock()
	cached, ok := cache[name]
	cacheMu.RUnlock()

	if ok {
		var setting map[string]interface{}
		err := json.Unmarshal(cached, &setting)
		if err != nil {
			return nil, err
		}
		return setting, nil
	}

	row, err := getSettingByName(name)
	if err != nil {
		return nil, err
	}

	var setting map[string]interface{}
	err = json.Unmarshal([]byte(row.Setting), &setting)
	if err != nil {
		return nil, err
	}

	cacheMu.Lock()
	cache[name] = []byte(row.Setting)
	cacheMu.Unlock()

	return setting, nil
}

// 是否需要验证码
func (m Setting) IsNeedCode(email string) bool {
	system, _ := SettingModel.GetSetting(SettingSystemKey)
	need := "no"
	switch system["login_code"] {
	case "open":
		need = "need"
	case "close":
		need = "no"
	default:
		code, found := core.Cache.Get("code::" + email)
		if found && code == "need" {
			need = "need"
		}
	}
	return need == "need"
}

// 是否需要验证码
func (m Setting) IsNeedInvite() bool {
	system, _ := SettingModel.GetSetting(SettingSystemKey)
	return system["reg"] == "invite"
}

// Setting 设置操作
func (m Setting) SettingOperation(name string, newSettings map[string]interface{}, isUpdate bool) (map[string]interface{}, error) {
	if newSettings == nil {
		return m.GetSetting(name)
	}

	err := m.UpdateSetting(name, newSettings, isUpdate)
	if err != nil {
		return nil, err
	}

	return m.GetSetting(name)
}

// UpdateSetting 更新设置
func (m Setting) UpdateSetting(name string, newSettings map[string]interface{}, isUpdate bool) error {
	serializedSettings, err := json.Marshal(newSettings)
	if err != nil {
		return err
	}

	err = updateSettingByName(name, string(serializedSettings))
	if err != nil {
		return err
	}

	cacheMu.Lock()
	cache[name] = serializedSettings
	cacheMu.Unlock()

	return nil
}

// getSettingByName 获取设置
func getSettingByName(name string) (*Setting, error) {
	settingModel := &Setting{}
	if err := core.DB.Where("name=?", name).First(settingModel).Error; err != nil {
		return nil, err
	}
	return settingModel, nil
}

// updateSettingByName 更新设置
func updateSettingByName(name string, setting string) error {
	settingModel := Setting{}
	if err := core.DB.Where("name=?", name).First(&settingModel).Error; err != nil {
		// 如果获取失败，则创建一个新的设置
		settingModel = Setting{
			Name:    name,
			Setting: setting,
		}
		if err := core.DB.Create(&settingModel).Error; err != nil {
			return err
		}
	} else {
		settingModel.Setting = setting
		if err := core.DB.Save(&settingModel).Error; err != nil {
			return err
		}
	}
	return nil
}

// SettingFind 获取设置
func (m Setting) SettingFind(name string, keyName string, defaultVal string) (string, error) {
	array, err := m.GetSetting(name)
	if err != nil {
		return defaultVal, err
	}
	if val, ok := array[keyName]; ok {
		return val.(string), nil
	}
	return defaultVal, nil
}
