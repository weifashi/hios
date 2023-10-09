package command

import (
	"fmt"
	"hios/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCommand.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "查看日志",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\nCommitSHA: %s\n", config.Version, config.CommitSHA)
	},
}
