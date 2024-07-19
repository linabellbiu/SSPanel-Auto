package service

import "github.com/robfig/cron/v3"

// 定时任务,会阻塞进程
func cronTable(f func()) {
	var crontab = cron.New()

	task := func() {
		f()
	}

	// 定时任务
	// 每天的凌晨1秒执行
	spec := "1 0 0 * * *"

	// 添加定时任务,
	_, _ = crontab.AddFunc(spec, task)

	crontab.Start()
	defer crontab.Stop()

	select {}
}
