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
	"os"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	// Import to initialize task
	_ "hios/app/task"
	"hios/app/wsc"

	ginI18n "github.com/gin-contrib/i18n"
)

var rootCommand = &cobra.Command{
	Use:   "hios",
	Short: "启动服务",
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
		if config.CONF.System.Prefix == "" {
			config.CONF.System.Prefix = "xw_hios_"
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
		// 初始化工作目录
		common.Mkdir(config.WorkDir, 0777)
		common.Mkdir(config.WorkDir+"/logs", 0777)
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

			// 设置日志输出到文件
			file, _ := os.OpenFile("./"+config.WorkDir+"/logs/request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			defer file.Close()
			gin.DefaultWriter = file
			gin.DefaultErrorWriter = file

			// 设置路由
			routers := gin.Default()
			routers.Use(middleware.OperationLog())
			routers.Use(gzip.Gzip(gzip.DefaultCompression))
			routers.SetHTMLTemplate(t)
			routers.Use(i18n.GinI18nLocalize())
			routers.SetFuncMap(template.FuncMap{
				"Localize": ginI18n.GetMessage,
			})
			routers.Any("/*path", func(context *gin.Context) {
				router.Init(context)
			})
			//
			common.PrintSuccess("启动成功: http://localhost:" + config.CONF.System.Port)
			//
			routers.Run(fmt.Sprintf("%s:%s", config.CONF.System.Host, config.CONF.System.Port))
			//
		}
	},
}

func Execute() {
	godotenv.Load(".env")
	rootCommand.CompletionOptions.DisableDefaultCmd = true
	flags := rootCommand.Flags()
	if os.Getenv("HIOS_WORKDIR") != "" {
		flags.StringVar(&config.WorkDir, "workDir", os.Getenv("HIOS_WORKDIR"), "工作目录")
	}
	flags.StringVar(&config.CONF.System.Host, "host", os.Getenv("HIOS_HOST"), "主机名，默认：0.0.0.0")
	flags.StringVar(&config.CONF.System.Port, "port", os.Getenv("HIOS_PORT"), "端口号，默认：30376")
	flags.StringVar(&config.CONF.System.Mode, "mode", os.Getenv("HIOS_MODE"), "运行模式，可选：debug|test|release")
	flags.StringVar(&config.CONF.System.Cache, "cache", os.Getenv("HIOS_CACHE"), "数据缓存目录，默认：{RunDir}/.cache")
	flags.StringVar(&config.CONF.System.WssUrl, "wss", os.Getenv("HIOS_WSS"), "服务端生成的url")
	flags.StringVar(&config.CONF.Jwt.SecretKey, "secret_key", "base64:ONdadQs1W4pY3h3dzr1jUSPrqLdsJQ9tCBZnb7HIDtk=", "jwt密钥")
	flags.StringVar(&config.CONF.Redis.RedisUrl, "redis_url", "redis://localhost:56379", "RedisUrl")
	flags.StringVar(&config.CONF.System.Prefix, "prefix", os.Getenv("HIOS_MYSQL_Prefix"), "数据前缀")
	flags.StringVar(&config.CONF.System.Dsn, "dsn", fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("HIOS_MYSQL_USERNAME"),
		os.Getenv("HIOS_MYSQL_PASSWORD"),
		os.Getenv("HIOS_MYSQL_HOST"),
		os.Getenv("HIOS_MYSQL_PORT"),
		os.Getenv("HIOS_MYSQL_DBNAME"),
	), "数据来源名称，如：sqlite://{CacheDir}/database.db")
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
