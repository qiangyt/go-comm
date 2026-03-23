package qshell

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/qiangyt/go-comm/v3/qlang"
)

// CommandOptions 创建命令的选项
type CommandOptionsT struct {
	// Ctx 上下文
	Ctx context.Context

	// Command 命令名
	Command string

	// Args 命令参数
	Args []string

	// MaxBufferSize 输出缓冲区最大大小（字节），默认 1MB
	MaxBufferSize int

	// Logger 日志记录器（可选）
	Logger qlang.Logger
}

type CommandOptions = *CommandOptionsT

// RunningCommand 正在运行的命令
type RunningCommandT struct {
	cmd        *exec.Cmd
	output     []byte
	mu         sync.Mutex
	done       chan struct{}
	maxBuffer  int
	logger     qlang.Logger
	terminalID string
	exitCode   *int
}

type RunningCommand = *RunningCommandT

// NewRunningCommand 创建并启动一个命令
func NewRunningCommand(options CommandOptions) RunningCommand {
	// 安全检查
	qlang.CheckTerminalCommand(options.Command, options.Args)

	ctx := options.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	maxBuffer := options.MaxBufferSize
	if maxBuffer <= 0 {
		maxBuffer = 1024 * 1024 // 默认 1MB
	}

	terminalID := fmt.Sprintf("terminal-%d-%d", time.Now().UnixNano(), time.Now().Nanosecond())

	cmd := exec.CommandContext(ctx, options.Command, options.Args...)

	// 创建管道捕获输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Errorf("创建 stdout 管道失败: %w", err))
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		panic(fmt.Errorf("创建 stderr 管道失败: %w", err))
	}

	result := &RunningCommandT{
		cmd:        cmd,
		done:       make(chan struct{}),
		maxBuffer:  maxBuffer,
		logger:     options.Logger,
		terminalID: terminalID,
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		panic(fmt.Errorf("启动命令失败: %w", err))
	}

	if result.logger != nil {
		result.logger.Info().Str("terminalId", terminalID).Int("pid", result.Pid()).Msg("终端命令已启动")
	}

	// 收集输出
	go func() {
		defer func() {
			qlang.RecoverAndLog(recover(), result.logger, "terminal output collector")
		}()
		defer close(result.done)

		// 合并 stdout 和 stderr
		reader := io.MultiReader(stdoutPipe, stderrPipe)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := append(scanner.Bytes(), '\n')
			result.mu.Lock()
			if len(result.output)+len(line) > result.maxBuffer {
				if len(result.output) > result.maxBuffer/2 {
					result.output = result.output[:result.maxBuffer/2]
				}
			}
			result.output = append(result.output, line...)
			result.mu.Unlock()
		}

		if result.logger != nil {
			result.logger.Info().Str("terminalId", terminalID).Msg("终端命令输出收集完成")
		}
	}()

	return result
}

// Pid 返回进程 ID
func (me RunningCommand) Pid() int {
	if me.cmd.Process == nil {
		return 0
	}
	return me.cmd.Process.Pid
}

// TerminalID 返回终端 ID
func (me RunningCommand) TerminalID() string {
	return me.terminalID
}

// GetOutput 获取当前输出（副本）
func (me RunningCommand) GetOutput() []byte {
	me.mu.Lock()
	output := make([]byte, len(me.output))
	copy(output, me.output)
	me.mu.Unlock()
	return output
}

// GetOutputString 获取当前输出的字符串形式
func (me RunningCommand) GetOutputString() string {
	return string(me.GetOutput())
}

// Wait 等待命令完成并返回退出码
func (me RunningCommand) Wait() int {
	<-me.done

	// 使用缓存避免重复调用 cmd.Wait()
	me.mu.Lock()
	if me.exitCode != nil {
		code := *me.exitCode
		me.mu.Unlock()
		return code
	}
	me.mu.Unlock()

	err := me.cmd.Wait()

	var code int
	if err != nil {
		if me.logger != nil {
			me.logger.Error(err).Str("terminalId", me.terminalID).Msg("终端命令执行失败")
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code = exitErr.ExitCode()
		} else {
			code = 1
		}
	} else {
		if me.logger != nil {
			me.logger.Info().Str("terminalId", me.terminalID).Msg("终端命令执行完成")
		}
		code = 0
	}

	me.mu.Lock()
	me.exitCode = &code
	me.mu.Unlock()

	return code
}

// Kill 终止命令
func (me RunningCommand) Kill() {
	if me.cmd.Process != nil {
		if err := me.cmd.Process.Kill(); err != nil {
			if me.logger != nil {
				me.logger.Error(err).Str("terminalId", me.terminalID).Msg("终止终端进程失败")
			}
		}
	} else {
		if me.logger != nil {
			me.logger.Warn().Str("terminalId", me.terminalID).Msg("终端进程不存在，无法终止")
		}
	}
}

// Done 返回完成信号 channel
func (me RunningCommand) Done() <-chan struct{} {
	return me.done
}

// IsDone 检查命令是否已完成
func (me RunningCommand) IsDone() bool {
	select {
	case <-me.done:
		return true
	default:
		return false
	}
}
