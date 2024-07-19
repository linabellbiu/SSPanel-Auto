package cmd

import (
	"fmt"
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

var (
	host   string
	email  string
	passwd string
	proxy  string
)

func RootCmd() error {
	rootCmd.AddCommand(CheckinCmd())
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&host, "host", "H", "", "需要访问的域名 例如: http://wwww.xxx.com (必填)")
	rootCmd.Flags().StringVarP(&email, "email", "e", "", "账户名,注册的邮箱账号 (必填)")
	rootCmd.Flags().StringVarP(&passwd, "passwd", "p", "", "密码,注册的密码 (必填)")
	rootCmd.Flags().StringVarP(&proxy, "proxy", "x", "", "设置http代理 例如:http://127.0.0.1:7890")

	_ = rootCmd.MarkFlagRequired("host")
	_ = rootCmd.MarkFlagRequired("email")
	_ = rootCmd.MarkFlagRequired("passwd")
}
