package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
