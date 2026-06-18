package qgin

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

// readRequestBody 读取 request body 并保留供后续处理
// 读取后会恢复 body，使得后续处理可以再次读取
func readRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// 读取 body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// 恢复 body 供后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}
