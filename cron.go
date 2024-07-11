package main

import (
	"github.com/robfig/cron/v3"
)

// 精确到秒
var crontab = cron.New(cron.WithSeconds())
