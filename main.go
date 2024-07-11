package main

// Import resty into your code and refer it as `resty`.
import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
)

type Option struct {
	Host       string `short:"d" long:"domain" description:"需要访问的域名 例如: http://wwww.xxx.com" required:"true"`
	Email      string `short:"e" long:"email" description:"账户名,注册的邮箱账号" required:"true"`
	Passwd     string `short:"p" long:"passwd" description:"密码,注册的密码" required:"true"`
	Remember   string `short:"r" long:"remember" description:"登录的请求参数" default:"week"`
	TryCount   int    `short:"n" long:"tryCount" description:"请求失败重试次数" default:"3"`
	HttpProxy  string `short:"t" long:"httpProxy" description:"设置http代理 例如:http://127.0.0.1:7890"`
	HttpsProxy string `short:"s" long:"httpsProxy" description:"设置https代理 例如:https://127.0.0.1:7890"`
}

var (
	cmd    Option
	client *resty.Client
)

func main() {
	_, err := flags.Parse(&cmd)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	client = resty.New()
	if cmd.HttpProxy != "" {
		_ = os.Setenv("http_proxy", cmd.HttpProxy)
	}
	if cmd.HttpsProxy != "" {
		_ = os.Setenv("http_proxy", cmd.HttpsProxy)
	}

	cookies, err := login()
	if err != nil {
		fmt.Println(fmt.Errorf("login failed %s", err.Error()))
		return
	}
	// 签到
	checkin(cookies)
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
