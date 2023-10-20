package model

import (
	"fmt"
	"hios/core"
	"os"
)

// User 用户实例
type User struct {
	Id          int         `gorm:"primary_key" json:"id"`
	Pinyin      string      `gorm:"type:varchar(255);comment:拼音（主要用于搜索）" json:"pinyin"`
	Email       string      `gorm:"type:varchar(100);uniqueIndex;comment:邮箱" json:"email"`
	Tel         string      `gorm:"type:varchar(50);comment:联系电话" json:"tel"`
	Nickname    string      `gorm:"type:varchar(255);comment:昵称" json:"nickname"`
	Profession  string      `gorm:"type:varchar(255);comment:职位/职称" json:"profession"`
	Avatar      string      `gorm:"type:varchar(255);comment:头像" json:"avatar"`
	Encrypt     string      `gorm:"type:varchar(50);comment:密钥" json:"encrypt,omitempty"`
	Password    string      `gorm:"type:varchar(50);comment:登录密码" json:"password,omitempty"`
	Changepass  int         `gorm:"type:tinyint(4);default:0;comment:登录需要修改密码" json:"changepass"`
	EmailVerity int         `gorm:"type:tinyint(1);default:0;comment:邮箱是否已验证" json:"email_verity"`
	Token       string      `gorm:"type:varchar(100);comment:token" json:"token,omitempty"`
	LoginNum    int         `gorm:"type:int(11);comment:累计登录次数" json:"login_num"`
	LastIp      string      `gorm:"type:varchar(20);comment:最后登录IP" json:"last_ip"`
	LastAt      core.TsTime `gorm:"comment:最后登录时间" json:"last_at"`
	LineIp      string      `gorm:"type:varchar(20);comment:最后在线IP（接口）" json:"line_ip"`
	LineAt      core.TsTime `gorm:"comment:最后在线时间（接口）" json:"line_at"`
	CreatedIp   string      `gorm:"type:varchar(20);comment:注册IP" json:"created_ip"`
	DisableAt   core.TsTime `gorm:"comment:禁用时间（离职时间）" json:"disable_at"`
	core.BaseAtModels
}

var (
	UserModel      = User{}
	UserBasicField = []string{"id", "email", "name", "avatar", "pinyin", "last_at", "disable_at", "created_at"}
)

// GetUserByID 根据用户ID获取用户信息
func (m User) GetUserByID(Userid int, filterSensitiveFields ...bool) (*User, error) {
	err := core.DB.Where("userid = ?", Userid).Find(&m).Error
	// 过滤掉密码等敏感字段
	if len(filterSensitiveFields) == 0 {
		m.FilterSensitiveFields()
	}
	return &m, err
}

// UpdatedEmailVerityByID 根据用户ID更新邮箱为已验证
func (m User) UpdatedEmailVerityByID(Userid int) {
	core.DB.Model(&m).Where("userid = ?", Userid).Update("email_verity", 1)
}

// UpdateUserByID 更新用户信息
func (m User) UpdateUserByID(userID int, data map[string]interface{}) (*User, error) {
	// 更新用户信息
	if err := core.DB.Model(&m).Where("userid = ?", userID).Updates(data).Error; err != nil {
		return nil, err
	}
	// 获取用户信息
	if err := core.DB.First(&m).Error; err != nil {
		return nil, err
	}
	m.FilterSensitiveFields()
	return &m, nil
}

// FilterSensitiveFields 过滤掉密码等敏感字段
func (m *User) FilterSensitiveFields() {
	m.Password = ""
	m.Encrypt = ""
}

// CheckSystem 检查环境是否允许
// onlyUserid 仅指定会员
func (m User) CheckSystem(onlyUserid ...int) bool {
	if len(onlyUserid) > 0 && onlyUserid[0] != m.Id {
		return true
	}
	if os.Getenv("PASSWORD_ADMIN") == "disabled" {
		if m.Id == 1 {
			return false
		}
	}
	if os.Getenv("PASSWORD_OWNER") == "disabled" {
		return false
	}
	return true
}

// GetOnlineStatus 获取在线状态
func (m User) GetOnlineStatus() bool {
	_, found := core.Cache.Get("User::online:" + fmt.Sprint(m.Id))
	if found {
		return true
	}
	err := core.DB.Where("userid = ?", m.Id).First(&WebSocket{}).Error
	return err == nil
}
