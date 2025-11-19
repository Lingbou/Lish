package commands

import (
	"context"
	"fmt"
	"os"
)

type CatCommand struct {
	stdout *os.File
}

func NewCatCommand(stdout *os.File) *CatCommand {
	return &CatCommand{stdout: stdout}
}

func (c *CatCommand) Name() string {
	return "cat"
}

func (c *CatCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定至少一个文件")
	}
	
	for _, filename := range args {
		content, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %w", filename, err)
		}
		
		fmt.Fprint(c.stdout, string(content))
	}
	
	return nil
}

func (c *CatCommand) Help() string {
	return `cat - 显示文件内容

用法:
  cat [文件...]

描述:
  连接文件并在标准输出上显示内容。

示例:
  cat file.txt           # 显示文件内容
  cat file1.txt file2.txt # 显示多个文件内容`
}

func (c *CatCommand) ShortHelp() string {
	return "显示文件内容"
}

