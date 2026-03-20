package benchmark

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/qiangyt/go-comm/v2"
	"github.com/qiangyt/go-comm/v2/qshell"
)

// ==================== 测试数据 ====================

var (
	smallJson       = `{"name": "test", "value": 123}`                                                                                          // 31 字节
	mediumJson      = `{"name": "medium_test", "items": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], "nested": {"level1": {"level2": {"level3": "deep"}}}}` // 128 字节
	largeJson       string
	smallJsonBytes  = 31
	mediumJsonBytes = 128
	largeJsonBytes  int
)

func init() {
	// 创建大型 JSON
	largeJson = `{"items": [`
	for i := 0; i < 100; i++ {
		if i > 0 {
			largeJson += ","
		}
		largeJson += fmt.Sprintf(`{"id": %d, "name": "item_%d", "value": %d}`, i, i, i*100)
	}
	largeJson += `]}`
	largeJsonBytes = len(largeJson)
}

// ==================== 标准库 encoding/json ====================

// BenchmarkStdlib_Unmarshal_small 使用标准库 encoding/json 解析小型 JSON (31 bytes)
func BenchmarkStdlib_Unmarshal_small(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(smallJson), &result)
	}
}

// BenchmarkStdlib_Unmarshal_medium 使用标准库 encoding/json 解析中型 JSON (128 bytes)
func BenchmarkStdlib_Unmarshal_medium(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(mediumJson), &result)
	}
}

// BenchmarkStdlib_Unmarshal_large 使用标准库 encoding/json 解析大型 JSON
func BenchmarkStdlib_Unmarshal_large(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		json.Unmarshal([]byte(largeJson), &result)
	}
}

// BenchmarkStdlib_Marshal_small 使用标准库 encoding/json 序列化小型对象
func BenchmarkStdlib_Marshal_small(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	for i := 0; i < b.N; i++ {
		json.Marshal(data)
	}
}

// BenchmarkStdlib_Marshal_medium 使用标准库 encoding/json 序列化中型对象
func BenchmarkStdlib_Marshal_medium(b *testing.B) {
	data := map[string]any{
		"name":  "medium_test",
		"items": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"nested": map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": "deep",
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		json.Marshal(data)
	}
}

// ==================== sonic 单线程 ====================

// BenchmarkSonic_Unmarshal_small 使用 sonic 解析小型 JSON (31 bytes)
func BenchmarkSonic_Unmarshal_small(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal([]byte(smallJson), &result)
	}
}

// BenchmarkSonic_Unmarshal_medium 使用 sonic 解析中型 JSON (128 bytes)
func BenchmarkSonic_Unmarshal_medium(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal([]byte(mediumJson), &result)
	}
}

// BenchmarkSonic_Unmarshal_large 使用 sonic 解析大型 JSON
func BenchmarkSonic_Unmarshal_large(b *testing.B) {
	var result map[string]any
	for i := 0; i < b.N; i++ {
		sonic.Unmarshal([]byte(largeJson), &result)
	}
}

// BenchmarkSonic_Marshal_small 使用 sonic 序列化小型对象
func BenchmarkSonic_Marshal_small(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	for i := 0; i < b.N; i++ {
		sonic.Marshal(data)
	}
}

// BenchmarkSonic_Marshal_medium 使用 sonic 序列化中型对象
func BenchmarkSonic_Marshal_medium(b *testing.B) {
	data := map[string]any{
		"name":  "medium_test",
		"items": []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"nested": map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": "deep",
				},
			},
		},
	}
	for i := 0; i < b.N; i++ {
		sonic.Marshal(data)
	}
}

// ==================== encoding/json 并发测试 ====================

// BenchmarkStdlib_Unmarshal_concurrent2 使用 2 个 goroutine 并发解析
func BenchmarkStdlib_Unmarshal_concurrent2(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				json.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Unmarshal_concurrent4 使用 4 个 goroutine 并发解析
func BenchmarkStdlib_Unmarshal_concurrent4(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				json.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Unmarshal_concurrent8 使用 8 个 goroutine 并发解析
func BenchmarkStdlib_Unmarshal_concurrent8(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 8; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				json.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Unmarshal_concurrent16 使用 16 个 goroutine 并发解析
func BenchmarkStdlib_Unmarshal_concurrent16(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 16; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				json.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Marshal_concurrent2 使用 2 个 goroutine 并发序列化
func BenchmarkStdlib_Marshal_concurrent2(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				json.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Marshal_concurrent4 使用 4 个 goroutine 并发序列化
func BenchmarkStdlib_Marshal_concurrent4(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				json.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Marshal_concurrent8 使用 8 个 goroutine 并发序列化
func BenchmarkStdlib_Marshal_concurrent8(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 8; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				json.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkStdlib_Marshal_concurrent16 使用 16 个 goroutine 并发序列化
func BenchmarkStdlib_Marshal_concurrent16(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 16; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				json.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// ==================== sonic 并发测试 ====================

// BenchmarkSonic_Unmarshal_concurrent2 使用 2 个 goroutine 并发解析
func BenchmarkSonic_Unmarshal_concurrent2(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				sonic.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Unmarshal_concurrent4 使用 4 个 goroutine 并发解析
func BenchmarkSonic_Unmarshal_concurrent4(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				sonic.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Unmarshal_concurrent8 使用 8 个 goroutine 并发解析
func BenchmarkSonic_Unmarshal_concurrent8(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 8; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				sonic.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Unmarshal_concurrent16 使用 16 个 goroutine 并发解析
func BenchmarkSonic_Unmarshal_concurrent16(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 16; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var result map[string]any
				sonic.Unmarshal([]byte(smallJson), &result)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Marshal_concurrent2 使用 2 个 goroutine 并发序列化
func BenchmarkSonic_Marshal_concurrent2(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sonic.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Marshal_concurrent4 使用 4 个 goroutine 并发序列化
func BenchmarkSonic_Marshal_concurrent4(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sonic.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Marshal_concurrent8 使用 8 个 goroutine 并发序列化
func BenchmarkSonic_Marshal_concurrent8(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 8; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sonic.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// BenchmarkSonic_Marshal_concurrent16 使用 16 个 goroutine 并发序列化
func BenchmarkSonic_Marshal_concurrent16(b *testing.B) {
	data := map[string]any{
		"name":  "test",
		"value": 123,
	}
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		for j := 0; j < 16; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sonic.Marshal(data)
			}()
		}
		wg.Wait()
	}
}

// ==================== go-comm 函数测试（动态配置对比） ====================

// BenchmarkComm_FromJson_Stdlib 使用 go-comm 的 FromJson 函数（使用 stdlib）
func BenchmarkComm_FromJson_Stdlib(b *testing.B) {
	// 保存当前配置
	currentBackend := comm.JSONConfig.Backend
	// 切换到 stdlib
	comm.JSONConfig.Backend = comm.JSONBackendStdlib
	var result map[string]any
	for i := 0; i < b.N; i++ {
		comm.FromJson(smallJson, false, &result)
	}
	// 恢复配置
	comm.JSONConfig.Backend = currentBackend
}

// BenchmarkComm_FromJson_Sonic 使用 go-comm 的 FromJson 函数（使用 sonic）
func BenchmarkComm_FromJson_Sonic(b *testing.B) {
	// 保存当前配置
	currentBackend := comm.JSONConfig.Backend
	// 切换到 sonic
	comm.JSONConfig.Backend = comm.JSONBackendSonic
	var result map[string]any
	for i := 0; i < b.N; i++ {
		comm.FromJson(smallJson, false, &result)
	}
	// 恢复配置
	comm.JSONConfig.Backend = currentBackend
}

// BenchmarkComm_ParseCommandOutput_Stdlib 使用 go-comm 的 ParseCommandOutput 函数（使用 stdlib）
func BenchmarkComm_ParseCommandOutput_Stdlib(b *testing.B) {
	// 保存当前配置
	currentBackend := comm.JSONConfig.Backend
	// 切换到 stdlib
	comm.JSONConfig.Backend = comm.JSONBackendStdlib
	jsonOutput := "$json$\n\n" + smallJson
	for i := 0; i < b.N; i++ {
		qshell.ParseCommandOutput(jsonOutput)
	}
	// 恢复配置
	comm.JSONConfig.Backend = currentBackend
}

// BenchmarkComm_ParseCommandOutput_Sonic 使用 go-comm 的 ParseCommandOutput 函数（使用 sonic）
func BenchmarkComm_ParseCommandOutput_Sonic(b *testing.B) {
	// 保存当前配置
	currentBackend := comm.JSONConfig.Backend
	// 切换到 sonic
	comm.JSONConfig.Backend = comm.JSONBackendSonic
	jsonOutput := "$json$\n\n" + smallJson
	for i := 0; i < b.N; i++ {
		qshell.ParseCommandOutput(jsonOutput)
	}
	// 恢复配置
	comm.JSONConfig.Backend = currentBackend
}
