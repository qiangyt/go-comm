package main

import "github.com/fastgh/go-comm"

func main() {
	comm.RunGoShellCommand(map[string]any{}, "", "zenity --info Hello")
}
