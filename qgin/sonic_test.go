package qgin

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/require"
)

// TestQginPackageImport 测试 qgin 包可以正常导入和使用
func TestQginPackageImport(t *testing.T) {
	// 测试可以安全地多次调用 ConfigureGinWithSonic
	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
	})

	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
	})
}

// TestConfigureGinWithSonicConfig 测试自定义配置
func TestConfigureGinWithSonicConfig(t *testing.T) {
	require.NotPanics(t, func() {
		ConfigureGinWithSonicConfig(sonic.Config{
			EscapeHTML: false,
		})
	})

	// 恢复默认配置
	ConfigureGinWithSonic()
}
