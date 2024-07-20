package cmd

import (
	"fmt"
	"github.com/linabellbiu/SSPanel-AutoCheckin/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "SSPanel-AutoCheckin",
	Short: "实现SSPanel框架搭建的平台流量自动签到功能,自动邀请注册的功能",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use \"SSPanel-AutoCheckin [command] --help\" for more information about a command.")
	},
	TraverseChildren: true,
}

func RootCmd() error {
	rootCmd.AddCommand(CheckinCmd())
	rootCmd.AddCommand(RegCmd())
	return rootCmd.Execute()
}

var (
	commonFlag = new(service.CommonFlag)
)

func init() {
	rootCmd.Flags().StringVarP(&commonFlag.Host, "host", "H", "", "需要访问的域名 例如: http://wwww.xxx.com (必填)")
	rootCmd.Flags().StringVarP(&commonFlag.Proxy, "proxy", "x", "", "设置http代理 例如:http://127.0.0.1:7890")

	_ = rootCmd.MarkFlagRequired("host")
}
