package qio

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v2/qerr"
)

// ReadBytesP ...
func ReadBytesP(reader io.Reader) []byte {
	r, err := ReadBytes(reader)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// ReadBytes ...
func ReadBytes(reader io.Reader) ([]byte, error) {
	r, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "read from Reader")
	}
	return r, nil
}

// ReadText ...
func ReadTextP(reader io.Reader) string {
	r, err := ReadText(reader)
	if err != nil {
		panic(qerr.NewSystemError(err.Error(), err))
	}
	return r
}

// ReadText ...
func ReadText(reader io.Reader) (string, error) {
	byts, err := io.ReadAll(reader)
	if err != nil {
		return "", errors.Wrap(err, "")
	}
	return string(byts), nil
}

func ReadLines(reader io.Reader) []string {
	r := make([]string, 0, 32)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		r = append(r, line)
	}

	return r
}

// CloseQuietly 关闭 io.Closer，忽略错误
// 用于资源可能已经关闭的场景，nil closer 不会 panic
func CloseQuietly(c io.Closer) {
	if c != nil {
		c.Close()
	}
}
