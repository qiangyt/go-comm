package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigureGinWithSonic(t *testing.T) {
	// 测试可以安全地多次调用
	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
	})

	require.NotPanics(t, func() {
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
		ConfigureGinWithSonic()
	})
}
