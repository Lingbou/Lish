package commands

import (
	"context"
	"fmt"
	"os"
)

type ClearCommand struct {
	stdout *os.File
}

func NewClearCommand(stdout *os.File) *ClearCommand {
	return &ClearCommand{stdout: stdout}
}

func (c *ClearCommand) Name() string {
	return "clear"
}

func (c *ClearCommand) Execute(ctx context.Context, args []string) error {
	// 使用 ANSI 转义码清屏
	fmt.Fprint(c.stdout, "\033[2J\033[H")
	return nil
}

func (c *ClearCommand) Help() string {
	return `clear - 清除终端屏幕

用法:
  clear

描述:
  清除终端屏幕内容。

示例:
  clear    # 清屏`
}

func (c *ClearCommand) ShortHelp() string {
	return "清除屏幕"
}
