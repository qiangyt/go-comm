package comm

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"mvdan.cc/sh/v3/interp"
)

// ============================================================
// TestHandlerRegistry 基本功能测试
// ============================================================

func TestHandlerRegistry_Register(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	a.NotNil(registry)

	// 注册处理器
	called := false
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = true
		fmt.Fprint(hc.Stdout, strings.Join(args, " "))
		return nil
	})

	// 验证注册成功
	handler, ok := registry.Match("echo")
	a.True(ok)
	a.NotNil(handler)

	// 执行处理器
	var out bytes.Buffer
	err := handler(context.Background(), interp.HandlerContext{Stdout: &out}, []string{"hello", "world"})
	a.NoError(err)
	a.True(called)
	a.Equal("hello world", out.String())
}

func TestHandlerRegistry_MatchNotFound(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()

	// 未注册的命令
	handler, ok := registry.Match("notexist")
	a.False(ok)
	a.Nil(handler)
}

func TestHandlerRegistry_MultipleHandlers(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "echo: "+strings.Join(args, " "))
		return nil
	})
	registry.Register("cat", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "cat: "+strings.Join(args, " "))
		return nil
	})

	// 测试 echo
	var out bytes.Buffer
	handler, ok := registry.Match("echo")
	a.True(ok)
	err := handler(context.Background(), interp.HandlerContext{Stdout: &out}, []string{"hello"})
	a.NoError(err)
	a.Equal("echo: hello", out.String())

	// 测试 cat
	out.Reset()
	handler, ok = registry.Match("cat")
	a.True(ok)
	err = handler(context.Background(), interp.HandlerContext{Stdout: &out}, []string{"file.txt"})
	a.NoError(err)
	a.Equal("cat: file.txt", out.String())
}

// ============================================================
// TestHandlerRegistry 匹配模式测试
// ============================================================

func TestHandlerRegistry_GlobMatch(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("git*", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "git handler")
		return nil
	}, WithMatchMode(MatchGlob))

	// git 应该匹配
	handler, ok := registry.Match("git")
	a.True(ok)
	a.NotNil(handler)

	// gitconfig 应该匹配
	handler, ok = registry.Match("gitconfig")
	a.True(ok)
	a.NotNil(handler)

	// go 不应该匹配
	handler, ok = registry.Match("go")
	a.False(ok)
	a.Nil(handler)
}

func TestHandlerRegistry_ExactMatch(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		return nil
	}, WithMatchMode(MatchExact))

	// echo 应该精确匹配
	handler, ok := registry.Match("echo")
	a.True(ok)

	// echotest 不应该匹配
	handler, ok = registry.Match("echotest")
	a.False(ok)
	_ = handler
}

func TestHandlerRegistry_RegexMatch(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register(`^git(-\w+)*$`, func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "matched")
		return nil
	}, WithMatchMode(MatchRegex))

	// git 应该匹配
	handler, ok := registry.Match("git")
	a.True(ok)

	// git-config 应该匹配
	handler, ok = registry.Match("git-config")
	a.True(ok)

	// gitconfig 不应该匹配（没有连字符）
	handler, ok = registry.Match("gitconfig")
	a.False(ok)
	_ = handler
}

// ============================================================
// TestHandlerRegistry 优先级测试
// ============================================================

func TestHandlerRegistry_Priority(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()

	// 注册低优先级处理器（默认优先级 0）
	registry.Register("test", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "low priority")
		return nil
	})

	// 注册高优先级处理器
	registry.Register("test", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		fmt.Fprint(hc.Stdout, "high priority")
		return nil
	}, WithPriority(100))

	// 应该匹配高优先级的处理器
	handler, ok := registry.Match("test")
	a.True(ok)

	var out bytes.Buffer
	err := handler(context.Background(), interp.HandlerContext{Stdout: &out}, nil)
	a.NoError(err)
	a.Equal("high priority", out.String())
}

func TestHandlerRegistry_PriorityOrder(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()

	// 注册多个不同优先级的处理器
	called := ""
	registry.Register("cmd", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = "priority-10"
		return nil
	}, WithPriority(10))

	registry.Register("cmd", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = "priority-100"
		return nil
	}, WithPriority(100))

	registry.Register("cmd", func(ctx context.Context, hc interp.HandlerContext, args []string) error {
		called = "priority-1"
		return nil
	}, WithPriority(1))

	// 应该匹配最高优先级的处理器
	handler, ok := registry.Match("cmd")
	a.True(ok)
	err := handler(context.Background(), interp.HandlerContext{}, nil)
	a.NoError(err)
	a.Equal("priority-100", called)
}

// ============================================================
// TestHandlerRegistry 其他功能测试
// ============================================================

func TestHandlerRegistry_List(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })
	registry.Register("cat", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })

	list := registry.List()
	a.Len(list, 2)
	a.Contains(list, "echo")
	a.Contains(list, "cat")
}

func TestHandlerRegistry_Count(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	a.Equal(0, registry.Count())

	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })
	a.Equal(1, registry.Count())

	registry.Register("cat", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })
	a.Equal(2, registry.Count())
}

func TestHandlerRegistry_Clear(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })
	a.Equal(1, registry.Count())

	registry.Clear()
	a.Equal(0, registry.Count())
}

// ============================================================
// TestHandlerOptions 选项测试
// ============================================================

func TestWithPriority(t *testing.T) {
	a := require.New(t)

	opt := WithPriority(50)
	a.NotNil(opt)

	entry := &handlerEntry{}
	opt(entry)
	a.Equal(50, entry.priority)
}

func TestWithMatchMode(t *testing.T) {
	a := require.New(t)

	opt := WithMatchMode(MatchRegex)
	a.NotNil(opt)

	entry := &handlerEntry{}
	opt(entry)
	a.Equal(MatchRegex, entry.matchMode)
}

func TestWithDescription(t *testing.T) {
	a := require.New(t)

	opt := WithDescription("test handler")
	a.NotNil(opt)

	entry := &handlerEntry{}
	opt(entry)
	a.Equal("test handler", entry.description)
}

// ============================================================
// TestHandlerRegistry Builder 模式测试
// ============================================================

func TestHandlerRegistry_Builder(t *testing.T) {
	a := require.New(t)

	// 测试链式调用
	registry := NewHandlerRegistry().
		Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil }).
		Register("cat", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })

	handler, ok := registry.Match("echo")
	a.True(ok)
	a.NotNil(handler)

	handler, ok = registry.Match("cat")
	a.True(ok)
	a.NotNil(handler)
}

// ============================================================
// TestMatchHandler 便捷函数测试
// ============================================================

func TestMatchHandler(t *testing.T) {
	a := require.New(t)

	registry := NewHandlerRegistry()
	registry.Register("echo", func(ctx context.Context, hc interp.HandlerContext, args []string) error { return nil })

	// 测试匹配
	handler, ok := MatchHandler(registry, "echo")
	a.True(ok)
	a.NotNil(handler)

	// 测试不匹配
	handler, ok = MatchHandler(registry, "notexist")
	a.False(ok)
	a.Nil(handler)

	// 测试 nil registry
	handler, ok = MatchHandler(nil, "echo")
	a.False(ok)
	a.Nil(handler)
}
