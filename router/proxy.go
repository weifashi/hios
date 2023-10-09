package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 生成代理
func CreatedProxy(c *gin.Context, urls string) {
	// 解析目标地址
	targetUrl, err := url.Parse(urls)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	// 创建反向代理处理程序
	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	// 修改请求头中的 Host 字段
	directorFunc := func(req *http.Request) {
		req.Host = targetUrl.Host
		req.URL.Scheme = targetUrl.Scheme
		req.URL.Host = targetUrl.Host
	}
	proxy.Director = directorFunc
	// 执行反向代理请求
	proxy.ServeHTTP(c.Writer, c.Request)
}

// 加载静态
func LoadStatic(c *gin.Context, urlPath string, appPath string, appPathIcon string) {
	if (strings.Contains(urlPath, "/assets/") || verifyIconPath(urlPath, appPathIcon)) && path.Ext(urlPath) != "" {
		if verifyIconPath(urlPath, appPathIcon) {
			c.File(appPathIcon + "/" + filepath.Base(urlPath))
		} else {
			c.File(appPath + "/dist/assets/" + filepath.Base(urlPath))
		}
	} else {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.File(appPath + "/dist/index.html")
	}
}

// 验证icon
func verifyIconPath(urlPath string, appPathIcon string) bool {
	return strings.Contains(urlPath, strings.Replace(appPathIcon, "plugins", "", -1))
}
