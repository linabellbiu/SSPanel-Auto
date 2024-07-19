package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "执行签到功能",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(host)
		fmt.Println(email)
		fmt.Println(passwd)
	},
}

func CheckinCmd() *cobra.Command {
	return checkinCmd
}

var (
	cronSpec    string
	cronDisable bool
	tryCount    int
)

func init() {
	checkinCmd.Flags().StringVarP(&cronSpec, "cron", "", "1 0 0 * * *", "设置每天定时执行,只用在本地执行,如果是放在github action中要关闭这个选项.\n配合'cron_disable'指令打开此功能")
	checkinCmd.Flags().BoolVarP(&cronDisable, "cron_disable", "", true, "关闭次指令后,可以使用'cron'设置定时执行")
	checkinCmd.Flags().IntVarP(&tryCount, "count", "", 3, "请求失败重试次数")
}
