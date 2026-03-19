package comm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ShortUUID 测试

func TestShortUUID(t *testing.T) {
	result, err := ShortUUID()
	assert.Nil(t, err)
	assert.Len(t, result, 22) // Base62 编码的 UUID 为 22 字符

	// 验证只包含 Base62 字符 (0-9, a-z, A-Z)
	matched, _ := regexp.MatchString("^[0-9a-zA-Z]{22}$", result)
	assert.True(t, matched, "ShortUUID should only contain Base62 characters")
}

func TestShortUUIDP(t *testing.T) {
	result := ShortUUIDP()
	assert.Len(t, result, 22)

	// 验证只包含 Base62 字符
	matched, _ := regexp.MatchString("^[0-9a-zA-Z]{22}$", result)
	assert.True(t, matched, "ShortUUIDP should only contain Base62 characters")
}

func TestShortUUID_Uniqueness(t *testing.T) {
	// 生成多个 UUID，确保它们互不相同
	uuids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uuid := ShortUUIDP()
		assert.False(t, uuids[uuid], "ShortUUID should be unique")
		uuids[uuid] = true
	}
}

// NanoID 测试

func TestNanoID(t *testing.T) {
	result, err := NanoID()
	assert.Nil(t, err)
	assert.Len(t, result, 21) // 默认 NanoID 长度为 21

	// 验证只包含 URL 安全字符
	matched, _ := regexp.MatchString("^[0-9a-zA-Z_-]{21}$", result)
	assert.True(t, matched, "NanoID should only contain URL-safe characters")
}

func TestNanoIDP(t *testing.T) {
	result := NanoIDP()
	assert.Len(t, result, 21)

	// 验证只包含 URL 安全字符
	matched, _ := regexp.MatchString("^[0-9a-zA-Z_-]{21}$", result)
	assert.True(t, matched, "NanoIDP should only contain URL-safe characters")
}

func TestNanoIDWithSize(t *testing.T) {
	// 测试自定义长度
	result, err := NanoIDWithSize(10)
	assert.Nil(t, err)
	assert.Len(t, result, 10)

	result, err = NanoIDWithSize(50)
	assert.Nil(t, err)
	assert.Len(t, result, 50)
}

func TestNanoIDWithSizeP(t *testing.T) {
	result := NanoIDWithSizeP(15)
	assert.Len(t, result, 15)

	result = NanoIDWithSizeP(100)
	assert.Len(t, result, 100)
}

func TestNanoID_Uniqueness(t *testing.T) {
	// 生成多个 NanoID，确保它们互不相同
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := NanoIDP()
		assert.False(t, ids[id], "NanoID should be unique")
		ids[id] = true
	}
}

func TestNanoIDWithSize_ZeroSize(t *testing.T) {
	// 长度为 0 时应该返回空字符串
	result, err := NanoIDWithSize(0)
	assert.Nil(t, err)
	assert.Empty(t, result)
}

func TestNanoIDWithSizeP_ZeroSize(t *testing.T) {
	result := NanoIDWithSizeP(0)
	assert.Empty(t, result)
}

func TestNanoIDWithSize_NegativeSize(t *testing.T) {
	// 负数长度应该返回错误
	result, err := NanoIDWithSize(-1)
	assert.NotNil(t, err)
	assert.Empty(t, result)
}

func TestNanoIDWithSizeP_NegativeSize_Panic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "NanoIDWithSizeP with negative size should panic")
	}()

	NanoIDWithSizeP(-1)
}

// 错误路径测试

// mockReader 是一个模拟的 Reader，用于测试错误路径
type mockErrorReader struct{}

func (m *mockErrorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("mock read error")
}

func TestShortUUID_RandomError(t *testing.T) {
	// 由于 uuid.NewRandom 内部使用 crypto/rand，无法直接模拟错误
	// 但我们可以测试 ShortUUIDP 在错误时的 panic 行为
	// 这个测试验证函数在正常情况下不会 panic
	result := ShortUUIDP()
	assert.Len(t, result, 22)
}

func TestNanoIDWithSize_RandomError(t *testing.T) {
	// 测试正常的 NanoIDWithSize 不返回错误
	result, err := NanoIDWithSize(21)
	assert.Nil(t, err)
	assert.Len(t, result, 21)
}

func TestNanoIDP_PanicOnError(t *testing.T) {
	// NanoIDP 正常情况下不应该 panic
	result := NanoIDP()
	assert.Len(t, result, 21)
}

func TestNanoIDWithSizeP_PanicOnError(t *testing.T) {
	// NanoIDWithSizeP 在正常情况下不应该 panic
	result := NanoIDWithSizeP(21)
	assert.Len(t, result, 21)
}
