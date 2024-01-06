package command

import (
	"fmt"
	"hios/app/service"
	"hios/config"
	"hios/core"
	"hios/database"
	"hios/utils/common"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// go run main.go url --type=node --uid=1
func init() {
	rootCommand.AddCommand(sysCmd)
	sysCmd.Flags().StringVarP(&types, "type", "t", "node", "Your type")
	sysCmd.Flags().StringVarP(&uid, "uid", "u", "", "Your uid")
}

var types string
var uid string
var sysCmd = &cobra.Command{
	Use:   "url",
	Short: "生成url",
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.CONF.System.Host == "" {
			config.CONF.System.Host = "0.0.0.0"
		}
		if config.CONF.System.Port == "" {
			config.CONF.System.Port = "30376"
		}
		if config.CONF.System.Cache == "" {
			config.CONF.System.Cache = common.RunDir("/.cache")
		}
		if config.CONF.System.Dsn == "mysql://:@tcp(:)/?charset=utf8mb4&parseTime=True&loc=Local" {
			config.CONF.System.Dsn = fmt.Sprintf("sqlite3://%s/%s", config.CONF.System.Cache, "database.db")
		}
		config.CONF.System.Start = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
		//
		err := common.WriteFile(config.CONF.System.Cache+"/config.json", common.StructToJson(config.CONF.System))
		if err != nil {
			common.PrintError(fmt.Sprintf("配置文件写入失败: %s", err.Error()))
			os.Exit(1)
		}
		// 初始化db
		err = core.InitDB()
		if err != nil {
			common.PrintError(fmt.Sprintf("数据库加载失败: %s", err.Error()))
			os.Exit(1)
		}
		// 初始化数据库
		err = database.Init()
		if err != nil {
			common.PrintError(fmt.Sprintf("数据库初始化失败: %s", err.Error()))
			os.Exit(1)
		}
	},
	Run: func(c *cobra.Command, args []string) {
		result := service.ClientService.CreateUrl(types, uid)
		fmt.Print(result)
	},
}
