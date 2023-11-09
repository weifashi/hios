package helper

import (
	"hios/i18n"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var ApiResponse = response{}

type response struct{}

// Response 返回结果
func (r response) Response(c *gin.Context, code int, msg string, values ...any) {
	c.Header("Expires", "-1")
	c.Header("Cache-Control", "no-cache")
	c.Header("Pragma", "no-cache")
	var data any
	if len(values) == 1 {
		data = values[0]
	} else if len(values) == 0 {
		data = gin.H{}
	} else {
		data = values
	}

	//
	if strings.Contains(c.GetHeader("Accept"), "application/json") ||
		strings.Contains(c.GetHeader("Content-Type"), "application/json") ||
		strings.Contains(c.GetHeader("X-Requested-With"), "XMLHttpRequest") ||
		strings.HasPrefix(c.Request.URL.Path, "/api/") ||
		!strings.Contains(c.Request.Method, http.MethodGet) {
		// 接口返回
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  msg,
			"data": data,
		})
	} else {
		// 页面返回
		if code == http.StatusMovedPermanently {
			c.Redirect(code, msg)
		} else {
			c.HTML(http.StatusOK, "index", gin.H{
				"CODE": code,
				"MSG":  url.QueryEscape(msg),
			})
		}
	}
}

// Success 成功
func (r response) Success(c *gin.Context, values ...any) {
	sgDetail := i18n.GetMsgWithMap("成功", map[string]any{"detail": values})
	r.Response(c, http.StatusOK, sgDetail, values...)
}

// Error 失败
func (r response) Error(c *gin.Context, msg string, values ...any) {
	msgDetail := i18n.GetMsgWithMap(msg, map[string]any{"detail": values})
	r.Response(c, http.StatusBadRequest, msgDetail)
}
