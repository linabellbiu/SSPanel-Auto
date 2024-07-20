package main

import "github.com/linabellbiu/SSPanel-Auto/cmd"

func main() {
	err := cmd.RootCmd()
	if err != nil {
		return
	}
}
