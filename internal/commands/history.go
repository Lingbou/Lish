package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

type HistoryCommand struct {
	stdout *os.File
}

func NewHistoryCommand(stdout *os.File) *HistoryCommand {
	return &HistoryCommand{stdout: stdout}
}

func (c *HistoryCommand) Name() string {
	return "history"
}

func (c *HistoryCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("history", pflag.ContinueOnError)
	clear := flags.BoolP("clear", "c", false, "清空历史记录")
	count := flags.IntP("count", "n", 0, "显示最近 N 条记录")

	if err := flags.Parse(args); err != nil {
		return err
	}

	historyFile := c.getHistoryFile()

	if *clear {
		return c.clearHistory(historyFile)
	}

	return c.showHistory(historyFile, *count)
}

func (c *HistoryCommand) getHistoryFile() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".lish_history"
	}
	return filepath.Join(homeDir, ".lish_history")
}

func (c *HistoryCommand) showHistory(historyFile string, count int) error {
	content, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintln(c.stdout, "历史记录为空")
			return nil
		}
		return fmt.Errorf("读取历史记录失败: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// 过滤空行
	var history []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			history = append(history, line)
		}
	}

	// 确定显示范围
	start := 0
	if count > 0 && count < len(history) {
		start = len(history) - count
	}

	// 显示历史记录
	for i := start; i < len(history); i++ {
		fmt.Fprintf(c.stdout, "%5d  %s\n", i+1, history[i])
	}

	return nil
}

func (c *HistoryCommand) clearHistory(historyFile string) error {
	if err := os.Remove(historyFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清空历史记录失败: %w", err)
	}

	fmt.Fprintln(c.stdout, "历史记录已清空")
	return nil
}

func (c *HistoryCommand) Help() string {
	return `history - 历史记录管理

用法:
  history                # 显示所有历史记录
  history -n 20          # 显示最近 20 条记录
  history -c             # 清空历史记录

选项:
  -n, --count   显示最近 N 条记录
  -c, --clear   清空历史记录

描述:
  管理命令历史记录。历史记录保存在 ~/.lish_history 文件中。

示例:
  history                # 查看所有历史
  history -n 10          # 查看最近 10 条
  history -c             # 清空历史`
}

func (c *HistoryCommand) ShortHelp() string {
	return "历史记录管理"
}
