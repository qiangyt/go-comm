package main

import (
	"github.com/qiangyt/go-comm/v2"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	f, _ := comm.CreateLockFile(fs, "/tmp/hi.pid", nil)
	f.Close()
}
