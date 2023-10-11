package migrations

import (
	"hios/app/model"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// WebSocket连接记录
func (s InitDatabase) AddTableWebSocket() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocket{})
		},
	}
}

// WebSocket临时消息
func (s InitDatabase) AddTableWebSocketTmpMsg() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketTmpMsg{})
		},
	}
}
