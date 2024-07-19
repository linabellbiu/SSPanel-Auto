package service

import (
	"encoding/json"
	"fmt"
	"os"
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
