package main

import (
	"fmt"

	"github.com/qiangyt/go-comm/v2"
)

func main() {
	// comm.RunGoshCommand(map[string]any{}, "", "zenity --info Hello", nil)
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("%+v", x)
		}
	}()
	comm.RunShellCommandP(map[string]string{}, "", "bash", "sudo echo hi",
		func() string { return "changeit" })
}
