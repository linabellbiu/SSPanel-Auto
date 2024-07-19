package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

type CheckinService struct {
	CronSpec    string
	CronDisable bool
	TryCount    int
	client      *resty.Client
	commonFlag  *CommonFlag
}

func (c *CheckinService) Run(commonFlag *CommonFlag) error {
	c.client = resty.New()
	c.commonFlag = commonFlag

	SetProxy(commonFlag.Proxy)

	cookies, err := login(c.client, c.commonFlag.Host, c.commonFlag.Email, c.commonFlag.Passwd, c.TryCount)
	if err != nil {
		return err
	}
	c.checkin(cookies)
	return nil
}

func (c *CheckinService) checkin(cookies []*http.Cookie) {
	for i := 0; i < c.TryCount; i++ {
		resp, err := c.client.R().SetCookies(cookies).Post(strings.TrimRight(c.commonFlag.Host, "/") + "/user/checkin")
		if err != nil {
			fmt.Println(fmt.Errorf("checkin fail err: %v", err.Error()))
			// 重新签到
			time.Sleep(3 * time.Second)
			continue
		}

		if resp.StatusCode() == http.StatusFound {
			fmt.Println(fmt.Printf("checkin http code %d,try login again", resp.StatusCode()))
			// 重新登录
			//if c.Login == "auto" {
			cookies, err = login(c.client, c.commonFlag.Host, c.commonFlag.Email, c.commonFlag.Passwd, c.TryCount)
			if err != nil {
				fmt.Println(fmt.Errorf("login failed %s", err.Error()))
				return
				//}
			}
			continue
		}

		if resp.StatusCode() != http.StatusOK {
			fmt.Println(fmt.Printf("Checkin http code %d,try checkin again", resp.StatusCode()))
			time.Sleep(3 * time.Second)
			continue
		}

		if err := Response(resp.Body()); err != nil {
			fmt.Println(fmt.Errorf("checkin err: %v \n", err))
			return
		}

		fmt.Println("Checkin success!!!")
		return
	}
	return
}

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
