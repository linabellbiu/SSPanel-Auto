package cmd

import (
	"github.com/linabellbiu/SSPanel-AutoCheckin/service"
	"github.com/spf13/cobra"
	"log"
)

var regCmd = &cobra.Command{
	Use:   "reg",
	Short: "执行注册功能",
	Run: func(cmd *cobra.Command, args []string) {
		if err := regService.Run(commonFlag); err != nil {
			log.Fatal("执行注册功能失败 err:" + err.Error())
			return
		}
	},
}

var regService = new(service.RegService)

func RegCmd() *cobra.Command {
	regCmd.Flags().StringVarP(&regService.Code, "code", "", "", "邀请码,填写才能反流量 (必填)")
	regCmd.Flags().StringVarP(&regService.CronSpec, "cron", "", "1 0 0 * * *", "设置每天定时执行,只用在本地执行,如果是放在github action中要关闭这个选项.\n配合'cron_disable'指令打开此功能")
	regCmd.Flags().BoolVarP(&regService.CronDisable, "cron_disable", "", true, "关闭次指令后,可以使用'cron'设置定时执行")
	regCmd.Flags().IntVarP(&regService.TryCount, "try", "", 3, "请求失败重试次数")
	regCmd.Flags().IntVarP(&regService.RegCount, "count", "", 3, "每次注册的数量,不要设置太大")

	_ = regCmd.MarkFlagRequired("code")
	return regCmd
}
