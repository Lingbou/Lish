package commands

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type EchoCommand struct {
	stdout *os.File
}

func NewEchoCommand(stdout *os.File) *EchoCommand {
	return &EchoCommand{stdout: stdout}
}

func (c *EchoCommand) Name() string {
	return "echo"
}

func (c *EchoCommand) Execute(ctx context.Context, args []string) error {
	// 展开环境变量
	var expanded []string
	for _, arg := range args {
		expanded = append(expanded, os.ExpandEnv(arg))
	}

	output := strings.Join(expanded, " ")
	fmt.Fprintln(c.stdout, output)

	return nil
}

func (c *EchoCommand) Help() string {
	return `echo - 输出文本

用法:
  echo [文本...]

描述:
  在标准输出上显示一行文本。支持环境变量展开。

环境变量:
  使用 $VAR 或 ${VAR} 格式引用环境变量。

示例:
  echo Hello World        # 输出: Hello World
  echo $HOME              # 输出用户主目录路径
  echo "Current: $PWD"    # 输出当前目录`
}

func (c *EchoCommand) ShortHelp() string {
	return "输出文本"
}
