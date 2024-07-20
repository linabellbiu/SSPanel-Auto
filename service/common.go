package service

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"math/rand"
	"os"
	"time"
)

type CommonFlag struct {
	Host   string
	Email  string
	Passwd string
	Proxy  string
}

type RespBody struct {
	Ret int64  `json:"ret"`
	Msg string `json:"msg"`
}

func Response(b []byte) error {
	var body RespBody
	if err := json.Unmarshal(b, &body); err != nil {
		return fmt.Errorf("unmarshal err: %v \n", err)
	}
	if body.Ret != 1 {
		return fmt.Errorf("msg: %s \n", body.Msg)
	}
	return nil
}

func SetProxy(proxy string) {
	if proxy != "" {
		_ = os.Setenv("http_proxy", fmt.Sprintf("http://%s", proxy))
		_ = os.Setenv("https_proxy", fmt.Sprintf("http://%s", proxy))
	}
}

type Task interface {
	start() error
}

// 定时任务,会阻塞进程
func cronTask(t Task, spec string) {
	var crontab = cron.New()

	task := func() {
		if err := t.start(); err != nil {
			fmt.Println(fmt.Errorf("cron task err: %v \n", err))
		}
	}

	// 添加定时任务,
	_, _ = crontab.AddFunc(spec, task)

	crontab.Start()
	defer crontab.Stop()

	select {}
}

func RandInt64() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(888888888) + r.Int63n(888888888) + 88888888 + time.Now().Unix()
}

func RandMinMax(min, max int64) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int63n(max-min) + min
}
