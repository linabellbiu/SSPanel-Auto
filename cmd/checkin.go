package cmd

import (
	"github.com/linabellbiu/SSPanel-AutoCheckin/service"
	"github.com/spf13/cobra"
	"log"
)

var checkinCmd = &cobra.Command{
	Use:   "checkin",
	Short: "执行签到功能",
	Run: func(cmd *cobra.Command, args []string) {
		if err := checkinService.Run(commonFlag); err != nil {
			log.Fatal("执行签到功能失败 err:" + err.Error())
			return
		}
	},
}

var checkinService = new(service.CheckinService)

func CheckinCmd() *cobra.Command {
	checkinCmd.Flags().StringVarP(&checkinService.CronSpec, "cron", "", "1 0 0 * * *", "设置每天定时执行,只用在本地执行,如果是放在github action中要关闭这个选项.\n配合'cron_disable'指令打开此功能")
	checkinCmd.Flags().BoolVarP(&checkinService.CronDisable, "cron_disable", "", true, "关闭次指令后,可以使用'cron'设置定时执行")
	checkinCmd.Flags().IntVarP(&checkinService.TryCount, "count", "", 3, "请求失败重试次数")
	return checkinCmd
}
