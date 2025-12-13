package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type DateCommand struct {
	stdout *os.File
}

func NewDateCommand(stdout *os.File) *DateCommand {
	return &DateCommand{stdout: stdout}
}

func (c *DateCommand) Name() string {
	return "date"
}

func (c *DateCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("date", pflag.ContinueOnError)
	format := flags.StringP("format", "f", "", "自定义格式")
	iso := flags.Bool("iso", false, "ISO 8601 格式")
	rfc := flags.Bool("rfc", false, "RFC 3339 格式")

	if err := flags.Parse(args); err != nil {
		return err
	}

	now := time.Now()

	// 根据选项选择格式
	var output string
	if *iso {
		output = now.Format("2006-01-02T15:04:05")
	} else if *rfc {
		output = now.Format(time.RFC3339)
	} else if *format != "" {
		// 支持常见格式占位符
		output = c.formatTime(now, *format)
	} else {
		// 默认格式
		output = now.Format("2006-01-02 15:04:05 Monday")
	}

	fmt.Fprintln(c.stdout, output)
	return nil
}

func (c *DateCommand) formatTime(t time.Time, format string) string {
	// 简单的格式转换
	replacements := map[string]string{
		"%Y": "2006",
		"%m": "01",
		"%d": "02",
		"%H": "15",
		"%M": "04",
		"%S": "05",
	}

	result := format
	for placeholder, goFormat := range replacements {
		if idx := findString(result, placeholder); idx >= 0 {
			result = result[:idx] + goFormat + result[idx+len(placeholder):]
		}
	}

	return t.Format(result)
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func (c *DateCommand) Help() string {
	return `date - 显示或设置日期时间

用法:
  date [选项]

选项:
  --iso           ISO 8601 格式
  --rfc           RFC 3339 格式
  -f, --format    自定义格式（%Y %m %d %H %M %S）

描述:
  显示当前日期和时间。

示例:
  date                    # 默认格式
  date --iso              # 2025-11-20T15:30:45
  date --rfc              # 2025-11-20T15:30:45+08:00
  date -f "%Y-%m-%d"      # 2025-11-20`
}

func (c *DateCommand) ShortHelp() string {
	return "显示日期时间"
}
