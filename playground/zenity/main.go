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
	comm.RunGoshCommandP(map[string]any{}, "", "gosh echo '$json$\n\ntrue'", nil)
}
