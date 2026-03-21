package main

import (
	"github.com/qiangyt/go-comm/v3/qio"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	f, _ := qio.CreateLockFile(fs, "/tmp/hi.pid", nil)
	f.Close()
}
