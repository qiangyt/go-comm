package qgin

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// bodyCaptureWriter 用于捕获 response body 同时写入原始 ResponseWriter
type bodyCaptureWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

// newBodyCaptureWriter 创建一个新的 bodyCaptureWriter
func newBodyCaptureWriter(w gin.ResponseWriter) *bodyCaptureWriter {
	return &bodyCaptureWriter{
		ResponseWriter: w,
		buffer:         &bytes.Buffer{},
	}
}

// Write 实现 io.Writer 接口，同时写入 buffer 和原始 ResponseWriter
func (w *bodyCaptureWriter) Write(data []byte) (int, error) {
	// 写入 buffer 用于捕获
	w.buffer.Write(data)
	// 同时写入原始 ResponseWriter
	return w.ResponseWriter.Write(data)
}

// WriteString 实现 io.StringWriter 接口
func (w *bodyCaptureWriter) WriteString(s string) (int, error) {
	// 写入 buffer 用于捕获
	w.buffer.WriteString(s)
	// 同时写入原始 ResponseWriter
	return w.ResponseWriter.WriteString(s)
}

// CapturedBody 返回捕获的 body 字符串
func (w *bodyCaptureWriter) CapturedBody() string {
	return w.buffer.String()
}

// Bytes 返回捕获的 body 字节切片
func (w *bodyCaptureWriter) Bytes() []byte {
	return w.buffer.Bytes()
}
