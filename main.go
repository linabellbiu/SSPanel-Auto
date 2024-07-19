package main

import "github.com/linabellbiu/SSPanel-AutoCheckin/cmd"

func main() {
	err := cmd.RootCmd()
	if err != nil {
		return
	}
}
