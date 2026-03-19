package qgin

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// ==================== bodyCaptureWriter ====================

func TestBodyCaptureWriter_CaptureData(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 写入数据
	data := []byte("hello world")
	n, err := capture.Write(data)

	a.NoError(err)
	a.Equal(len(data), n)
	a.Equal("hello world", capture.CapturedBody())
}

func TestBodyCaptureWriter_WriteToOriginal(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 写入数据
	data := []byte("hello world")
	capture.Write(data)

	// 验证数据也写入到了原始 ResponseWriter
	a.Equal("hello world", w.Body.String())
}

func TestBodyCaptureWriter_WriteString(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 写入字符串
	n, err := capture.WriteString("hello world")

	a.NoError(err)
	a.Equal(11, n)
	a.Equal("hello world", capture.CapturedBody())
	a.Equal("hello world", w.Body.String())
}

func TestBodyCaptureWriter_MultipleWrites(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 多次写入
	capture.Write([]byte("hello "))
	capture.Write([]byte("world"))
	capture.WriteString("!")

	a.Equal("hello world!", capture.CapturedBody())
	a.Equal("hello world!", w.Body.String())
}

func TestBodyCaptureWriter_Empty(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 未写入任何数据
	a.Equal("", capture.CapturedBody())
	a.Equal("", w.Body.String())
}

func TestBodyCaptureWriter_Bytes(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 写入数据
	data := []byte("hello world")
	capture.Write(data)

	// 验证 Bytes() 方法
	a.Equal(data, capture.Bytes())
}

func TestBodyCaptureWriter_DelegateOtherMethods(t *testing.T) {
	a := require.New(t)

	// 创建 mock response writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 设置一些 header
	c.Writer.Header().Set("Content-Type", "application/json")

	// 创建 bodyCaptureWriter
	capture := newBodyCaptureWriter(c.Writer)

	// 验证 Header() 委托到原始 writer
	a.Equal("application/json", capture.Header().Get("Content-Type"))
}

// ==================== 比较测试：bytes.Buffer vs 直接 string 拼接 ====================

func TestBodyCaptureWriter_CompareBufferVsString(t *testing.T) {
	a := require.New(t)

	// 使用 bytes.Buffer
	buf := &bytes.Buffer{}
	buf.Write([]byte("hello "))
	buf.Write([]byte("world"))
	resultBuffer := buf.String()

	// 使用直接 string 拼接（应该使用 buffer，因为性能更好）
	// 这里只是展示 bytes.Buffer 的行为

	a.Equal("hello world", resultBuffer)
}
