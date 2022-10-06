package main

import (
	"fmt"

	"github.com/fastgh/go-comm/v2"
)

func main() {
	// comm.RunGoshCommand(map[string]any{}, "", "zenity --info Hello", nil)
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("%+v", x)
		}
	}()
	comm.RunShellCommandP(map[string]any{}, "", "bash", "sudo echo hi",
		func() string { return "changeit" })
}
