package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
			return nil, fmt.Errorf("login body unmarshal err: %v \n", err)
		}

		fmt.Printf("login success \n cookie: %v \n", resp.Cookies())
		return resp.Cookies(), nil
	}

	return nil, errors.New("login failed")
}
