package core

type BaseIdModels struct {
	Id int `gorm:"primary_key" json:"id"`
}

type BaseAtModels struct {
	CreatedAt TsTime `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt TsTime `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}
