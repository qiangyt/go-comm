package qcoll

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/qiangyt/go-comm/v2/qerr"
)

// HashCalculator 哈希计算器
type HashCalculatorT struct{}

type HashCalculator = *HashCalculatorT

// NewHashCalculator 创建哈希计算器
func NewHashCalculator() HashCalculator {
	return &HashCalculatorT{}
}

// CalculateMD5 计算文件的MD5哈希值
func (h HashCalculator) CalculateMD5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(qerr.NewSystemError("open file for MD5 calculation", err))
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(qerr.NewSystemError("calculate MD5", err))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// CalculateSHA256 计算文件的SHA256哈希值
func (h HashCalculator) CalculateSHA256(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(qerr.NewSystemError("open file for SHA256 calculation", err))
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(qerr.NewSystemError("calculate SHA256", err))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

// CalculateMD5FromReader 从io.Reader计算MD5
func (h HashCalculator) CalculateMD5FromReader(r io.Reader) string {
	hash := md5.New()
	if _, err := io.Copy(hash, r); err != nil {
		panic(qerr.NewSystemError("calculate MD5 from reader", err))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// CalculateSHA256FromReader 从io.Reader计算SHA256
func (h HashCalculator) CalculateSHA256FromReader(r io.Reader) string {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		panic(qerr.NewSystemError("calculate SHA256 from reader", err))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
