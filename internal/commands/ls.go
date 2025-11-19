package commands

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m" // 目录
	colorGreen  = "\033[32m" // 可执行文件
	colorCyan   = "\033[36m" // 符号链接
	colorYellow = "\033[33m" // 特殊文件
)

type LsCommand struct {
	stdout *os.File
}

func NewLsCommand(stdout *os.File) *LsCommand {
	return &LsCommand{stdout: stdout}
}

func (c *LsCommand) Name() string {
	return "ls"
}

func (c *LsCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("ls", pflag.ContinueOnError)
	longFormat := flags.BoolP("long", "l", false, "使用长格式")
	all := flags.BoolP("all", "a", false, "显示隐藏文件")
	humanReadable := flags.BoolP("human-readable", "h", false, "以人类可读的格式显示大小")

	if err := flags.Parse(args); err != nil {
		return err
	}

	// 获取要列出的目录
	dirs := flags.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	for i, dir := range dirs {
		if i > 0 {
			fmt.Fprintln(c.stdout)
		}

		if len(dirs) > 1 {
			fmt.Fprintf(c.stdout, "%s:\n", dir)
		}

		if err := c.listDir(dir, *longFormat, *all, *humanReadable); err != nil {
			return err
		}
	}

	return nil
}

func (c *LsCommand) listDir(dir string, longFormat, all, humanReadable bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 过滤隐藏文件
	var filtered []fs.DirEntry
	for _, entry := range entries {
		if !all && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filtered = append(filtered, entry)
	}

	// 排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name() < filtered[j].Name()
	})

	if longFormat {
		return c.printLongFormat(dir, filtered, humanReadable)
	}

	return c.printSimpleFormat(filtered)
}

func (c *LsCommand) printSimpleFormat(entries []fs.DirEntry) error {
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			fmt.Fprintf(c.stdout, "%s%s%s  ", colorBlue, name, colorReset)
		} else if isExecutable(entry) {
			fmt.Fprintf(c.stdout, "%s%s%s  ", colorGreen, name, colorReset)
		} else {
			fmt.Fprintf(c.stdout, "%s  ", name)
		}
	}
	fmt.Fprintln(c.stdout)
	return nil
}

func (c *LsCommand) printLongFormat(dir string, entries []fs.DirEntry, humanReadable bool) error {
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return err
		}

		// 文件类型和权限
		mode := info.Mode()
		modeStr := formatMode(mode)

		// 文件大小
		size := info.Size()
		sizeStr := formatSize(size, humanReadable)

		// 修改时间
		modTime := info.ModTime().Format("2006-01-02 15:04")

		// 文件名（带颜色）
		name := entry.Name()
		coloredName := name
		if entry.IsDir() {
			coloredName = colorBlue + name + colorReset
		} else if isExecutable(entry) {
			coloredName = colorGreen + name + colorReset
		}

		fmt.Fprintf(c.stdout, "%s %10s %s %s\n", modeStr, sizeStr, modTime, coloredName)
	}
	return nil
}

func formatMode(mode fs.FileMode) string {
	str := ""

	// 文件类型
	if mode.IsDir() {
		str += "d"
	} else if mode&fs.ModeSymlink != 0 {
		str += "l"
	} else {
		str += "-"
	}

	// 权限
	str += formatPerm(mode, 6) // owner
	str += formatPerm(mode, 3) // group
	str += formatPerm(mode, 0) // other

	return str
}

func formatPerm(mode fs.FileMode, shift uint) string {
	perm := (mode >> shift) & 7
	str := ""

	if perm&4 != 0 {
		str += "r"
	} else {
		str += "-"
	}

	if perm&2 != 0 {
		str += "w"
	} else {
		str += "-"
	}

	if perm&1 != 0 {
		str += "x"
	} else {
		str += "-"
	}

	return str
}

func formatSize(size int64, humanReadable bool) string {
	if !humanReadable {
		return fmt.Sprintf("%d", size)
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

	units := []string{"K", "M", "G", "T", "P"}
	return fmt.Sprintf("%.1f%s", float64(size)/float64(div), units[exp])
}

func isExecutable(entry fs.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

func (c *LsCommand) Help() string {
	return `ls - 列出目录内容

用法:
  ls [选项] [目录...]

选项:
  -l, --long            使用长格式显示详细信息
  -a, --all             显示所有文件，包括隐藏文件
  -h, --human-readable  以人类可读的格式显示文件大小（K, M, G）

描述:
  列出目录中的文件和子目录。默认显示当前目录的内容。

示例:
  ls              # 列出当前目录
  ls -l           # 长格式显示
  ls -la          # 显示所有文件，包括隐藏文件
  ls -lh /root    # 以人类可读格式显示 /root 目录`
}

func (c *LsCommand) ShortHelp() string {
	return "列出目录内容"
}
