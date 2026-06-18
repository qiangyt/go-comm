package qgin

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// errorReader 是一个模拟读取错误的 reader
type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (r *errorReader) Close() error {
	return nil
}

// ==================== readRequestBody ====================

func TestReadRequestBody_JsonBody(t *testing.T) {
	a := require.New(t)

	// 创建请求
	body := `{"name":"test","value":123}`
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 读取 request body
	result := readRequestBody(c)

	a.Equal(body, result)
	// 验证 body 可以被再次读取（应该被恢复）
	bodyBytes, err := io.ReadAll(c.Request.Body)
	a.NoError(err)
	a.Equal(body, string(bodyBytes))
}

func TestReadRequestBody_EmptyBody(t *testing.T) {
	a := require.New(t)

	// 创建空请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 读取 request body
	result := readRequestBody(c)

	a.Equal("", result)
}

func TestReadRequestBody_NilBody(t *testing.T) {
	a := require.New(t)

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	c.Request.Body = nil // 显式设置 body 为 nil

	// 读取 request body
	result := readRequestBody(c)

	a.Equal("", result)
}

func TestReadRequestBody_LargeBody(t *testing.T) {
	a := require.New(t)

	// 创建大请求体
	body := make([]byte, 10000)
	for i := range body {
		body[i] = 'a'
	}
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 读取 request body
	result := readRequestBody(c)

	a.Len(result, 10000)
}

func TestReadRequestBody_FormData(t *testing.T) {
	a := require.New(t)

	// 创建 form 请求
	body := "name=test&value=123"
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 读取 request body
	result := readRequestBody(c)

	a.Equal(body, result)
}

func TestReadRequestBody_BinaryBody(t *testing.T) {
	a := require.New(t)

	// 创建二进制请求体
	body := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/octet-stream")

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 读取 request body
	result := readRequestBody(c)

	// 二进制 body 应该被读取为字符串
	a.Equal(string(body), result)
}

func TestReadRequestBody_ReadError(t *testing.T) {
	a := require.New(t)

	// 创建 gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", &errorReader{})

	// 读取 request body，应该返回空字符串（错误处理）
	result := readRequestBody(c)

	a.Equal("", result)
}
