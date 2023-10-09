package v1

import (
	"hios/app/model"

	"github.com/gin-gonic/gin"
)

type BaseApi struct {
	Route    string       `mapstructure:"route"`
	Token    string       `mapstructure:"token"`
	Userinfo *model.User  `mapstructure:"userinfo"`
	Context  *gin.Context `mapstructure:"context"`
}

type NotAuthBaseApi struct {
	BaseApi BaseApi `mapstructure:"base_api"`
}
