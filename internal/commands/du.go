package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type DuCommand struct {
	stdout *os.File
}

func NewDuCommand(stdout *os.File) *DuCommand {
	return &DuCommand{stdout: stdout}
}

func (c *DuCommand) Name() string {
	return "du"
}

func (c *DuCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("du", pflag.ContinueOnError)
	humanReadable := flags.BoolP("human-readable", "h", false, "人类可读格式")
	summarize := flags.BoolP("summarize", "s", false, "只显示总计")

	if err := flags.Parse(args); err != nil {
		return err
	}

	paths := flags.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		size, err := c.calculateSize(path, !*summarize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "du: %v\n", err)
			continue
		}

		sizeStr := c.formatSize(size, *humanReadable)
		fmt.Fprintf(c.stdout, "%s\t%s\n", sizeStr, path)
	}

	return nil
}

func (c *DuCommand) calculateSize(path string, showSubdirs bool) (int64, error) {
	var totalSize int64

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略无法访问的文件
		}

		if !info.IsDir() {
			totalSize += info.Size()
		}

		return nil
	})

	return totalSize, err
}

func (c *DuCommand) formatSize(size int64, humanReadable bool) string {
	if !humanReadable {
		return fmt.Sprintf("%d", size/1024) // KB
	}

	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	units := []string{"K", "M", "G", "T"}
	return fmt.Sprintf("%.1f%s", float64(size)/float64(div), units[exp])
}

func (c *DuCommand) Help() string {
	return `du - 显示磁盘使用情况

用法:
  du [选项] [文件或目录...]

选项:
  -h, --human-readable  以人类可读格式显示（K, M, G）
  -s, --summarize       只显示总计

描述:
  估算文件和目录的磁盘使用量。

示例:
  du                    # 当前目录使用情况
  du -h                 # 人类可读格式
  du -sh .              # 只显示总计
  du -h dir1 dir2       # 多个目录`
}

func (c *DuCommand) ShortHelp() string {
	return "磁盘使用情况"
}
