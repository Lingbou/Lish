package commands

import (
	"context"
	"io"
)

// Command 定义命令接口
type Command interface {
	// Name 返回命令名称
	Name() string

	// Execute 执行命令
	Execute(ctx context.Context, args []string) error

	// Help 返回帮助信息
	Help() string

	// ShortHelp 返回简短描述
	ShortHelp() string
}

// Context 命令执行上下文
type Context struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	Env     map[string]string
	WorkDir string
}
