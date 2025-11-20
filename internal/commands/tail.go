package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type TailCommand struct {
	stdout *os.File
}

func NewTailCommand(stdout *os.File) *TailCommand {
	return &TailCommand{stdout: stdout}
}

func (c *TailCommand) Name() string {
	return "tail"
}

func (c *TailCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("tail", pflag.ContinueOnError)
	lines := flags.IntP("lines", "n", 10, "显示的行数")
	follow := flags.BoolP("follow", "f", false, "实时监控文件变化")
	
	if err := flags.Parse(args); err != nil {
		return err
	}
	
	files := flags.Args()
	if len(files) == 0 {
		return fmt.Errorf("tail: 需要指定文件")
	}
	
	if *follow {
		// 实时监控模式（只支持单个文件）
		if len(files) > 1 {
			return fmt.Errorf("tail -f: 只能监控一个文件")
		}
		return c.followFile(ctx, files[0], *lines)
	}
	
	// 普通模式
	for i, filename := range files {
		if i > 0 {
			fmt.Fprintln(c.stdout)
		}
		
		if len(files) > 1 {
			fmt.Fprintf(c.stdout, "==> %s <==\n", filename)
		}
		
		if err := c.printTail(filename, *lines); err != nil {
			return err
		}
	}
	
	return nil
}

func (c *TailCommand) printTail(filename string, n int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	
	// 读取所有行
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		return err
	}
	
	// 打印最后 N 行
	start := 0
	if len(lines) > n {
		start = len(lines) - n
	}
	
	for i := start; i < len(lines); i++ {
		fmt.Fprintln(c.stdout, lines[i])
	}
	
	return nil
}

func (c *TailCommand) followFile(ctx context.Context, filename string, n int) error {
	// 先打印最后 N 行
	if err := c.printTail(filename, n); err != nil {
		return err
	}
	
	// 打开文件准备监控
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	
	// 移动到文件末尾
	file.Seek(0, 2)
	
	// 监控文件变化
	scanner := bufio.NewScanner(file)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			// 尝试读取新内容
			for scanner.Scan() {
				fmt.Fprintln(c.stdout, scanner.Text())
			}
			
			if err := scanner.Err(); err != nil {
				return err
			}
		}
	}
}

func (c *TailCommand) Help() string {
	return `tail - 显示文件尾部内容

用法:
  tail [选项] 文件...

选项:
  -n, --lines   显示的行数（默认 10 行）
  -f, --follow  实时监控文件变化（类似 tail -f）

描述:
  显示文件的后 N 行内容。使用 -f 可以实时监控文件更新。

示例:
  tail file.txt           # 显示最后 10 行
  tail -n 20 file.txt     # 显示最后 20 行
  tail -f log.txt         # 实时监控日志文件（Ctrl+C 退出）`
}

func (c *TailCommand) ShortHelp() string {
	return "显示文件尾部"
}

