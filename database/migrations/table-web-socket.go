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

// WebSocket对话
func (s InitDatabase) AddTableWebSocketDialog() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketDialog{})
		},
	}
}

// 对话消息记录
func (s InitDatabase) AddTableWebSocketDialogMsg() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketDialogMsg{})
		},
	}
}

// 对话消息阅读记录
func (s InitDatabase) AddTableWebSocketDialogMsgRead() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketDialogMsgRead{})
		},
	}
}

// 对话消息待办记录
func (s InitDatabase) AddTableWebSocketDialogMsgTodo() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketDialogMsgTodo{})
		},
	}
}

// 对话用户
func (s InitDatabase) AddTableWebSocketDialogUser() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2023061500",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&model.WebSocketDialogUser{})
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
