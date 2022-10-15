package main

import (
	"github.com/fastgh/go-comm/v2"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	f, _ := comm.CreateLockFile(fs, "/tmp/hi.pid")
	f.Close()
}
