package command

import (
	"fmt"
	"hios/config"
	"hios/core"
	"hios/database"
	"hios/i18n"
	"hios/router"
	"hios/router/middleware"
	"hios/utils/common"
	"hios/web"
	"html/template"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	// Import to initialize task
	_ "hios/app/task"
	"hios/app/wsc"
)

var rootCommand = &cobra.Command{
	Use:   "hios",
	Short: "启动服务",
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.CONF.System.Host == "" {
			config.CONF.System.Host = "0.0.0.0"
		}
		if config.CONF.System.Port == "" {
			config.CONF.System.Port = "3376"
		}
		if config.CONF.System.Cache == "" {
			config.CONF.System.Cache = common.RunDir("/.cache")
		}
		if config.CONF.System.Dsn == "" {
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
		//
		common.PrintSuccess("启动成功: http://localhost:" + config.CONF.System.Port)
	},
	Run: func(cmd *cobra.Command, args []string) {

		// 连接服务端
		if config.CONF.System.WssUrl != "" {
			wsc.WorkStart()
		} else {
			// 启动服务端
			t, err := template.New("index").Parse(string(web.IndexByte))
			if err != nil {
				common.PrintError(fmt.Sprintf("模板解析失败: %s", err.Error()))
				os.Exit(1)
			}
			if config.CONF.System.Mode == "debug" {
				gin.SetMode(gin.DebugMode)
			} else if config.CONF.System.Mode == "test" {
				gin.SetMode(gin.TestMode)
			} else {
				gin.SetMode(gin.ReleaseMode)
			}
			gin.DefaultWriter = io.Discard
			r := gin.Default()
			r.Use(middleware.CustomLogger())
			r.Use(middleware.OperationLog())
			r.Use(gzip.Gzip(gzip.DefaultCompression))
			r.Use(i18n.GinI18nLocalize())
			r.SetFuncMap(template.FuncMap{
				"Localize": ginI18n.GetMessage,
			})
			r.SetHTMLTemplate(t)
			r.Any("/*path", func(context *gin.Context) {
				router.Init(context)
			})
			_ = r.Run(fmt.Sprintf("%s:%s", config.CONF.System.Host, config.CONF.System.Port))
		}

	},
}

func Execute() {
	godotenv.Load(".env")
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	rootCommand.Flags().StringVar(&config.CONF.System.Host, "host", os.Getenv("HOST"), "主机名，默认：0.0.0.0")
	rootCommand.Flags().StringVar(&config.CONF.System.Port, "port", os.Getenv("PORT"), "端口号，默认：3376")
	rootCommand.Flags().StringVar(&config.CONF.System.Mode, "mode", os.Getenv("MODE"), "运行模式，可选：debug|test|release")
	rootCommand.Flags().StringVar(&config.CONF.System.Cache, "cache", "", "数据缓存目录，默认：{RunDir}/.cache")
	rootCommand.Flags().StringVar(&config.CONF.System.WssUrl, "wss", "", "服务端生成的url")
	rootCommand.Flags().StringVar(&config.CONF.System.Dsn, "dsn", "", "数据来源名称，如：sqlite://{CacheDir}/database.db")
	rootCommand.Flags().StringVar(&config.CONF.Jwt.SecretKey, "secret_key", "base64:ONdadQs1W4pY3h3dzr1jUSPrqLdsJQ9tCBZnb7HIDtk=", "jwt密钥")
	rootCommand.Flags().StringVar(&config.CONF.Redis.RedisUrl, "redis_url", "redis://localhost:56379", "RedisUrl")
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
