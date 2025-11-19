package commands

import (
	"context"
	"fmt"
	"os"
)

type PwdCommand struct {
	stdout *os.File
}

func NewPwdCommand(stdout *os.File) *PwdCommand {
	return &PwdCommand{stdout: stdout}
}

func (c *PwdCommand) Name() string {
	return "pwd"
}

func (c *PwdCommand) Execute(ctx context.Context, args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}

	fmt.Fprintln(c.stdout, dir)
	return nil
}

func (c *PwdCommand) Help() string {
	return `pwd - 显示当前工作目录

用法:
  pwd

描述:
  打印当前工作目录的完整路径。

示例:
  pwd    # 显示当前目录`
}

func (c *PwdCommand) ShortHelp() string {
	return "显示当前工作目录"
}
