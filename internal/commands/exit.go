package commands

import (
	"context"
	"os"
	"strconv"
)

type ExitCommand struct {
	exitFunc func(int)
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{
		exitFunc: os.Exit,
	}
}

func (c *ExitCommand) Name() string {
	return "exit"
}

func (c *ExitCommand) Execute(ctx context.Context, args []string) error {
	exitCode := 0
	
	if len(args) > 0 {
		code, err := strconv.Atoi(args[0])
		if err == nil {
			exitCode = code
		}
	}
	
	c.exitFunc(exitCode)
	return nil
}

func (c *ExitCommand) Help() string {
	return `exit - 退出 Shell

用法:
  exit [退出码]

描述:
  退出 Lish Shell。可以指定退出码（默认为 0）。

参数:
  退出码  可选，指定退出状态码（默认 0）

示例:
  exit      # 正常退出
  exit 1    # 以错误状态退出`
}

func (c *ExitCommand) ShortHelp() string {
	return "退出 Shell"
}

