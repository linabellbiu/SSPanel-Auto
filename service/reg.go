package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

type RegService struct {
	CronSpec    string
	CronDisable bool
	TryCount    int
	RegCount    int
	Code        string
	client      *resty.Client
	commonFlag  *CommonFlag
}

func (c *RegService) Run(commonFlag *CommonFlag) error {
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

func (c *RegService) start() error {
	if c.CronDisable {
		c.reg()
	} else {
		cronTask(c, c.CronSpec)
	}
	return nil
}

func (c *RegService) reg() {
	for i := 0; i < c.RegCount; i++ {
		for i := 0; i < c.TryCount; i++ {
			reqBody := c.generateBody()
			resp, err := c.client.R().SetBody(reqBody).Post(strings.TrimRight(c.commonFlag.Host, "/") + "/auth/register")
			if err != nil {
				fmt.Println(fmt.Errorf("reg fail err: %v", err.Error()))
				time.Sleep(3 * time.Second)
				continue
			}

			if resp.StatusCode() != http.StatusOK {
				fmt.Println(fmt.Printf("reg http code %d,try checkin again", resp.StatusCode()))
				time.Sleep(3 * time.Second)
				continue
			}

			if err := Response(resp.Body()); err != nil {
				fmt.Println(fmt.Errorf("reg err: %v \n", err))
				break
			}
			fmt.Println("reg success!!!")
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}

func (c *RegService) generateBody() map[string]interface{} {
	// 随机生产账户名
	name := uuid.New().String()
	// 随机生成账户
	email := fmt.Sprintf("%d@gmail.com", RandInt64())
	// 随机生成密码
	pwd := fmt.Sprintf("%d", RandInt64())
	// 联络方式
	wechat := fmt.Sprintf("%d", RandInt64())
	// 联络方式类型
	imType := RandMinMax(1, 4)

	q := map[string]interface{}{
		"email":    email,
		"name":     name,
		"passwd":   pwd,
		"wechat":   wechat,
		"imtype":   imType,
		"repasswd": pwd,
		"code":     c.Code,
	}

	log.Printf("reg param: %+v \n", q)
	return q
}
