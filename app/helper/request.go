package helper

import (
	"fmt"
	"hios/app/constant"
	"hios/utils/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

var ApiRequest = request{}

type request struct{}

// Token 获取Token（Header、Query、Cookie）
func (r request) Token(c *gin.Context) string {
	token := c.GetHeader("token")
	if token == "" {
		token = r.Input(c, "token")
	}
	if token == "" {
		token = r.Cookie(c, "token")
	}
	return token
}

// Input 获取参数（优先POST、取Query）
func (r request) Input(c *gin.Context, key string) string {
	if c.PostForm(key) != "" {
		return strings.TrimSpace(c.PostForm(key))
	}
	return strings.TrimSpace(c.Query(key))
}

// Scheme 获取Scheme
func (r request) Scheme(c *gin.Context) string {
	scheme := "http://"
	if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https://"
	}
	return scheme
}

// Home 获取Home Url
func (r request) Home(c *gin.Context) string {
	return fmt.Sprintf("%s%s", r.Scheme(c), c.Request.Host)
}

// Cookie 获取Cookie
func (r request) Cookie(c *gin.Context, name string) string {
	value, _ := c.Cookie(name)
	return value
}

// 获取当前域名
func (r request) GetCurrentDomain(c *gin.Context) string {
	if r.IsHttps(c) {
		return fmt.Sprintf("https://%s", c.Request.Host)
	} else {
		return fmt.Sprintf("http://%s", c.Request.Host)
	}
}

// 判断当前请求是否https
func (r request) IsHttps(c *gin.Context) bool {
	return c.Request.TLS != nil
}

// SetCookie 设置Cookie
func (r request) SetCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetCookie(name, value, maxAge, "/", "", false, false)
}

// RemoveCookie 删除Cookie
func (r request) RemoveCookie(c *gin.Context, name string) {
	c.SetCookie(name, "", -1, "/", "", false, false)
}

// 绑定参数
func (r request) ShouldBindAll(c *gin.Context, obj any) {
	if err := c.ShouldBind(obj); err != nil {
		logger.Error(err.Error())
		ApiResponse.Error(c, constant.ErrInvalidParameter, err)
		panic(nil)
	}
	if err := c.ShouldBindJSON(obj); err != nil && err.Error() != "EOF" {
		logger.Error(err.Error())
		ApiResponse.Error(c, constant.ErrInvalidParameter, err)
		panic(nil)
	}
}
