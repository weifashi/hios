package command

import (
	"bufio"
	"fmt"
	"hios/config"
	"io"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "查看日志",
	Run: func(c *cobra.Command, args []string) {

		f := "./" + config.WorkDir + "/logs/request.log"

		if len(args) > 0 {
			if args[0] == "wsc" || args[0] == "c" {
				f = "./" + config.WorkDir + "/logs/wsc.log"
			}
			if args[0] == "wss" || args[0] == "s" {
				f = "./" + config.WorkDir + "/logs/wss.log"
			}
			if args[0] == "request" || args[0] == "r" {
				f = "./" + config.WorkDir + "/logs/request.log"
			}
		}

		//
		cmd := exec.Command("tail", "-f", f)

		// 创建一个管道，用于获取命令的标准输出
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("创建管道出错:", err)
			return
		}

		// 启动命令
		if err := cmd.Start(); err != nil {
			fmt.Println("启动命令出错:", err)
			return
		}

		// 读取命令的标准输出
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("读取输出出错:", err)
				continue
			}

			// 处理读取到的内容
			fmt.Print(line)
		}

		// 等待命令执行完成
		if err := cmd.Wait(); err != nil {
			fmt.Println("等待命令完成出错:", err)
			return
		}
	},
}
