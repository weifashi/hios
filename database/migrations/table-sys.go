package migrations

import (
	"hios/app/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// 添加设置表
func (s InitDatabase) AddTableSetting() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023060500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.Setting{})
		},
	}
}
