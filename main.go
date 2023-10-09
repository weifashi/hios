package main

import (
	"hios/command"
	"runtime"
)

// @title FileWarehouse
// @version 1.0
// @description  FileWarehouse是一款轻量级的开源在线项目文件管理工具，提供各类文档协作工具。
// @termsOfService https://file.weifashi.cn/
// @license.name AGPL-3.0 license
// @license.url http://www.gnu.org/licenses/
// @host http://localhost
// @BasePath /api/v1

//go:generate swag init --parseDependency -o ./app/docs -g ./main.go -d ./app -g ../main.go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	command.Execute()
}
