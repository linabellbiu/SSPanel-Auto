package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

type CheckinService struct {
	Email       string
	Passwd      string
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

	if c.CronDisable {
		return c.start()
	} else {
		cronTask(c, c.CronSpec)
	}

	return nil
}

func (c *CheckinService) start() error {
	cookies, err := login(c.client, c.commonFlag.Host, c.Email, c.Passwd, c.TryCount)
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
