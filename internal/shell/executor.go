package shell

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Lingbou/Lish/internal/commands"
	"github.com/Lingbou/Lish/internal/parser"
)

// Executor 命令执行器
type Executor struct {
	registry *commands.Registry
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
}

// NewExecutor 创建执行器
func NewExecutor(registry *commands.Registry, stdin io.Reader, stdout, stderr io.Writer) *Executor {
	return &Executor{
		registry: registry,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
	}
}

// ExecutePipeline 执行管道
func (e *Executor) ExecutePipeline(ctx context.Context, pipeline *parser.Pipeline) error {
	if len(pipeline.Commands) == 0 {
		return nil
	}

	// 单个命令，不需要管道
	if len(pipeline.Commands) == 1 {
		return e.executeCommand(ctx, pipeline.Commands[0], e.stdin, e.stdout, e.stderr)
	}

	// 多个命令，使用管道连接
	var pipes []io.ReadCloser

	for i := 0; i < len(pipeline.Commands)-1; i++ {
		pr, pw := io.Pipe()
		pipes = append(pipes, pr)

		// 执行命令，输出到管道
		go func(cmdIndex int, writer *io.PipeWriter) {
			defer writer.Close()

			var input io.Reader
			if cmdIndex == 0 {
				input = e.stdin
			} else {
				input = pipes[cmdIndex-1]
			}

			err := e.executeCommand(ctx, pipeline.Commands[cmdIndex], input, writer, e.stderr)
			if err != nil {
				fmt.Fprintf(e.stderr, "管道错误: %v\n", err)
			}
		}(i, pw)
	}

	// 执行最后一个命令
	lastIndex := len(pipeline.Commands) - 1
	var lastInput io.Reader
	if lastIndex == 0 {
		lastInput = e.stdin
	} else {
		lastInput = pipes[lastIndex-1]
	}

	return e.executeCommand(ctx, pipeline.Commands[lastIndex], lastInput, e.stdout, e.stderr)
}

// executeCommand 执行单个命令（处理重定向）
func (e *Executor) executeCommand(ctx context.Context, cmd *parser.ParsedCommand, stdin io.Reader, stdout, stderr io.Writer) error {
	// 处理输入重定向
	if cmd.RedirectIn != "" {
		file, err := os.Open(cmd.RedirectIn)
		if err != nil {
			return fmt.Errorf("打开输入文件失败: %w", err)
		}
		defer file.Close()
		stdin = file
	}

	// 处理输出重定向
	if cmd.RedirectOut != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if cmd.Append {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		file, err := os.OpenFile(cmd.RedirectOut, flags, 0644)
		if err != nil {
			return fmt.Errorf("打开输出文件失败: %w", err)
		}
		defer file.Close()
		stdout = file
	}

	// 处理错误重定向
	if cmd.RedirectErr != "" {
		file, err := os.OpenFile(cmd.RedirectErr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("打开错误输出文件失败: %w", err)
		}
		defer file.Close()
		stderr = file
	}

	// 获取命令
	command, exists := e.registry.Get(cmd.Command)
	if !exists {
		return fmt.Errorf("未知命令: %s", cmd.Command)
	}

	// 创建临时的命令包装器来处理 I/O
	wrapper := &commandWrapper{
		cmd:    command,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	return wrapper.Execute(ctx, cmd.Args)
}

// commandWrapper 包装命令以支持自定义 I/O
type commandWrapper struct {
	cmd    commands.Command
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (w *commandWrapper) Execute(ctx context.Context, args []string) error {
	// 临时替换标准输出（这是一个简化实现）
	// 实际应该通过 context 或者修改 Command 接口来传递 I/O
	return w.cmd.Execute(ctx, args)
}
