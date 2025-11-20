package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type HeadCommand struct {
	stdout *os.File
}

func NewHeadCommand(stdout *os.File) *HeadCommand {
	return &HeadCommand{stdout: stdout}
}

func (c *HeadCommand) Name() string {
	return "head"
}

func (c *HeadCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("head", pflag.ContinueOnError)
	lines := flags.IntP("lines", "n", 10, "显示的行数")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) == 0 {
		return fmt.Errorf("head: 需要指定文件")
	}

	for i, filename := range files {
		if i > 0 {
			fmt.Fprintln(c.stdout)
		}

		if len(files) > 1 {
			fmt.Fprintf(c.stdout, "==> %s <==\n", filename)
		}

		if err := c.printHead(filename, *lines); err != nil {
			return err
		}
	}

	return nil
}

func (c *HeadCommand) printHead(filename string, n int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() && count < n {
		fmt.Fprintln(c.stdout, scanner.Text())
		count++
	}

	return scanner.Err()
}

func (c *HeadCommand) Help() string {
	return `head - 显示文件头部内容

用法:
  head [选项] 文件...

选项:
  -n, --lines  显示的行数（默认 10 行）

描述:
  显示文件的前 N 行内容。

示例:
  head file.txt           # 显示前 10 行
  head -n 20 file.txt     # 显示前 20 行
  head -n 5 *.txt         # 显示多个文件的前 5 行`
}

func (c *HeadCommand) ShortHelp() string {
	return "显示文件头部"
}
