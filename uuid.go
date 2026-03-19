package comm

import (
	"crypto/rand"
	"fmt"
	"math/big"

	uuid "github.com/google/uuid"
)

// Base62 字符集
const base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NanoID URL 安全字符集
const nanoIDCharset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_-"

// ShortUUID 生成 Base62 编码的 UUID（22 字符）
// UUID 是 128 位（16 字节），Base62 编码后为 22 字符
func ShortUUID() (string, error) {
	// 生成新的 UUID
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("生成 UUID 失败: %w", err)
	}

	// 将 UUID 转换为 Base62
	return encodeBase62(id[:]), nil
}

// ShortUUIDP 生成 Base62 编码的 UUID（22 字符），失败时 panic
func ShortUUIDP() string {
	result, err := ShortUUID()
	if err != nil {
		panic(NewSystemError("生成 ShortUUID 失败", err))
	}
	return result
}

// encodeBase62 将字节数组编码为 Base62 字符串
func encodeBase62(data []byte) string {
	// 将字节转换为大整数
	num := new(big.Int).SetBytes(data)
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	result := make([]byte, 0, 22)
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		result = append(result, base62Charset[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	// 填充前导零以达到 22 字符
	for len(result) < 22 {
		result = append([]byte{'0'}, result...)
	}

	return string(result)
}

// NanoID 生成 21 字符的 NanoID（默认长度）
func NanoID() (string, error) {
	return NanoIDWithSize(21)
}

// NanoIDP 生成 21 字符的 NanoID（默认长度），失败时 panic
func NanoIDP() string {
	return NanoIDWithSizeP(21)
}

// NanoIDWithSize 生成指定长度的 NanoID
func NanoIDWithSize(size int) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("NanoID 长度不能为负数: %d", size)
	}
	if size == 0 {
		return "", nil
	}

	charsetLen := big.NewInt(int64(len(nanoIDCharset)))
	result := make([]byte, size)

	for i := 0; i < size; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("生成 NanoID 失败: %w", err)
		}
		result[i] = nanoIDCharset[n.Int64()]
	}

	return string(result), nil
}

// NanoIDWithSizeP 生成指定长度的 NanoID，失败时 panic
func NanoIDWithSizeP(size int) string {
	result, err := NanoIDWithSize(size)
	if err != nil {
		panic(NewSystemError("生成 NanoID 失败", err))
	}
	return result
}
