package service

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
	"time"
)

func login(client *resty.Client, host, email, passwd string, tryCount int) ([]*http.Cookie, error) {
	for i := 0; i < tryCount; i++ {
		resp, err := client.R().SetBody(map[string]interface{}{
			"email":       email,
			"passwd":      passwd,
			"remember_me": "week",
		}).Post(strings.TrimRight(host, "/") + "/auth/login")
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
