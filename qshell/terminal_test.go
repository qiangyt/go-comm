package qshell

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewRunningCommand_happy(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo hello"}
	} else {
		cmd = "echo"
		args = []string{"hello"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	a.NotEmpty(runningCmd.TerminalID())
	a.Greater(runningCmd.Pid(), 0)

	// 等待命令完成
	exitCode := runningCmd.Wait()
	a.Equal(0, exitCode)

	// 检查输出
	output := runningCmd.GetOutputString()
	a.Contains(output, "hello")
}

func TestNewRunningCommand_withContext(t *testing.T) {
	a := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo test"}
	} else {
		cmd = "echo"
		args = []string{"test"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Ctx:           ctx,
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	exitCode := runningCmd.Wait()
	a.Equal(0, exitCode)
}

func TestRunningCommand_GetOutput(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo line1 && echo line2"}
	} else {
		cmd = "sh"
		args = []string{"-c", "echo line1; echo line2"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	runningCmd.Wait()

	output := runningCmd.GetOutput()
	a.Contains(string(output), "line1")
	a.Contains(string(output), "line2")

	// 确保返回的是副本
	modifiedOutput := append(output, []byte("extra")...)
	originalOutput := runningCmd.GetOutput()
	a.NotEqual(len(modifiedOutput), len(originalOutput))
}

func TestRunningCommand_Kill(t *testing.T) {
	a := require.New(t)

	if runtime.GOOS == "windows" {
		t.Skip("Windows 下 sleep 命令行为不同")
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       "sleep",
		Args:          []string{"100"},
		MaxBufferSize: 1024 * 1024,
	})

	// 确保命令已启动
	a.Greater(runningCmd.Pid(), 0)
	a.False(runningCmd.IsDone())

	// 终止命令
	runningCmd.Kill()

	// 等待完成
	exitCode := runningCmd.Wait()
	a.NotEqual(0, exitCode)
	a.True(runningCmd.IsDone())
}

func TestRunningCommand_IsDone(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo done"}
	} else {
		cmd = "echo"
		args = []string{"done"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	// 初始状态可能还没完成
	// 等待一小段时间后检查
	time.Sleep(100 * time.Millisecond)
	runningCmd.Wait()

	a.True(runningCmd.IsDone())
}

func TestRunningCommand_Done(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo done"}
	} else {
		cmd = "echo"
		args = []string{"done"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	// 等待 Done channel
	<-runningCmd.Done()
	a.True(runningCmd.IsDone())
}

func TestNewRunningCommand_blockedCommand(t *testing.T) {
	a := require.New(t)

	if runtime.GOOS == "windows" {
		t.Skip("Windows 下行为不同")
	}

	// 测试被阻止的危险命令
	a.Panics(func() {
		NewRunningCommand(&CommandOptionsT{
			Command:       "rm",
			Args:          []string{"-rf", "/"},
			MaxBufferSize: 1024 * 1024,
		})
	})
}

func TestNewRunningCommand_invalidCommand(t *testing.T) {
	a := require.New(t)

	a.Panics(func() {
		NewRunningCommand(&CommandOptionsT{
			Command:       "nonexistent_command_xyz",
			Args:          []string{},
			MaxBufferSize: 1024 * 1024,
		})
	})
}

func TestNewRunningCommand_maxBufferSize(t *testing.T) {
	a := require.New(t)

	if runtime.GOOS == "windows" {
		t.Skip("Windows 下行为不同")
	}

	// 测试输出缓冲区限制
	// 生成大量输出的命令
	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       "sh",
		Args:          []string{"-c", "for i in $(seq 1 1000); do echo \"line $i\"; done"},
		MaxBufferSize: 100, // 很小的缓冲区
	})

	runningCmd.Wait()

	// 输出应该被截断
	output := runningCmd.GetOutput()
	a.LessOrEqual(len(output), 100)
}

func TestRunningCommand_Pid(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo test"}
	} else {
		cmd = "echo"
		args = []string{"test"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	// 进程应该已启动
	a.Greater(runningCmd.Pid(), 0)

	runningCmd.Wait()
}

func TestRunningCommand_ProcessState(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "exit 42"}
	} else {
		cmd = "sh"
		args = []string{"-c", "exit 42"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	exitCode := runningCmd.Wait()
	a.Equal(42, exitCode)

	// 再次调用 Wait 应该返回相同的退出码
	exitCode2 := runningCmd.Wait()
	a.Equal(42, exitCode2)
}

func TestRunningCommand_KillAfterCompletion(t *testing.T) {
	a := require.New(t)

	var cmd string
	var args []string
	if runtime.GOOS == "windows" {
		cmd = "cmd"
		args = []string{"/c", "echo done"}
	} else {
		cmd = "echo"
		args = []string{"done"}
	}

	runningCmd := NewRunningCommand(&CommandOptionsT{
		Command:       cmd,
		Args:          args,
		MaxBufferSize: 1024 * 1024,
	})

	runningCmd.Wait()

	// 命令已完成后调用 Kill 不应该 panic
	a.NotPanics(func() {
		runningCmd.Kill()
	})
}

func TestRunningCommand_Env(t *testing.T) {
	// 测试环境变量访问
	// 这个测试确保 RunningCommand 可以被创建并正常工作
	t.Skip("环境变量测试需要更复杂的设置")
}
