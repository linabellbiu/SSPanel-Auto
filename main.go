package main

// Import resty into your code and refer it as `resty`.
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Option struct {
	Host       string `short:"d" long:"domain" description:"需要访问的域名 例如: http://wwww.xxx.com" required:"true"`
	Email      string `short:"e" long:"email" description:"账户名,注册的邮箱账号" required:"true"`
	Passwd     string `short:"p" long:"passwd" description:"密码,注册的密码" required:"true"`
	Remember   string `short:"r" long:"remember" description:"登录的请求参数" default:"week"`
	TryCount   int    `short:"n" long:"tryCount" description:"请求失败重试次数" default:"3"`
	HttpProxy  string `short:"t" long:"httpProxy" description:"设置http代理 例如:http://127.0.0.1:7890"`
	HttpsProxy string `short:"s" long:"httpsProxy" description:"设置https代理 例如:https://127.0.0.1:7890"`
	Login      string `short:"l" long:"login" default:"auto" description:"设置身份是否过期自动登录 \n auto \n always \n no \n"`
	Cron       bool   `short:"C" long:"cron" description:"设置每天定时执行,只用在本地执行,如果是放在github action中要关闭这个选项. \n 每天的执行时间: 1 0 0 * * * "`
}

var (
	cmd    Option
	client = resty.New()
	// 精确到秒
	crontab = cron.New(cron.WithSeconds())
)

func main() {
	_, err := flags.Parse(&cmd)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	fmt.Println("flags:", cmd)

	if cmd.Login != "auto" && cmd.Login != "always" && cmd.Login != "no" {
		log.Fatal("请设置设置身份是否过期自动登录,查看详情 `./SSPanel-AutoCheckin -h`")
	}

	if cmd.HttpProxy != "" {
		_ = os.Setenv("http_proxy", cmd.HttpProxy)
	}
	if cmd.HttpsProxy != "" {
		_ = os.Setenv("http_proxy", cmd.HttpsProxy)
	}

	if cmd.Cron {
		cronTable(run)
	} else {
		run()
	}
}

// 定时任务,会阻塞进程
func cronTable(f func()) {
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

func run() {
	cookies, err := login()
	if err != nil {
		fmt.Println(fmt.Errorf("login failed %s", err.Error()))
		return
	}
	// 签到
	checkin(cookies)
}

func checkin(cookies []*http.Cookie) {
	for i := 0; i < cmd.TryCount; i++ {
		resp, err := client.R().SetCookies(cookies).Post(strings.TrimRight(cmd.Host, "/") + "/user/checkin")
		if err != nil {
			fmt.Println(fmt.Errorf("checkin fail err: %v", err.Error()))
			// 重新签到
			time.Sleep(3 * time.Second)
			continue
		}

		if resp.StatusCode() == http.StatusFound {
			fmt.Println(fmt.Printf("checkin http code %d,try login again", resp.StatusCode()))
			// 重新登录
			if cmd.Login == "auto" {
				cookies, err = login()
				if err != nil {
					fmt.Println(fmt.Errorf("login failed %s", err.Error()))
					return
				}
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

func login() ([]*http.Cookie, error) {
	for i := 0; i < cmd.TryCount; i++ {
		resp, err := client.R().SetBody(map[string]interface{}{
			"email":       cmd.Email,
			"passwd":      cmd.Passwd,
			"remember_me": cmd.Remember,
		}).Post(strings.TrimRight(cmd.Host, "/") + "/auth/login")
		if err != nil {
			fmt.Println(fmt.Errorf("login fail err: %v,try login again", err.Error()))
			// 重新登录
			time.Sleep(3 * time.Second)
			continue
		}

		if resp.StatusCode() != http.StatusOK {
			fmt.Println(fmt.Errorf("login fail,http code %d,try login again", resp.StatusCode()))
			// 重新登录
			time.Sleep(3 * time.Second)
			continue
		}

		if err := Response(resp.Body()); err != nil {
			return nil, fmt.Errorf("login body unmarshal err: %v ,body :%s\n", err, string(resp.Body()))
		}

		fmt.Printf("login success \n cookie: %v \n", resp.Cookies())
		return resp.Cookies(), nil
	}

	return nil, errors.New("login failed")
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
