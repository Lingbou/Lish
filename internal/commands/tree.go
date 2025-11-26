package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

type TreeCommand struct {
	stdout *os.File
}

func NewTreeCommand(stdout *os.File) *TreeCommand {
	return &TreeCommand{stdout: stdout}
}

func (c *TreeCommand) Name() string {
	return "tree"
}

func (c *TreeCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("tree", pflag.ContinueOnError)
	level := flags.IntP("level", "L", 0, "显示的层级深度（0=无限制）")
	dirOnly := flags.BoolP("directories", "d", false, "只显示目录")

	if err := flags.Parse(args); err != nil {
		return err
	}

	paths := flags.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		fmt.Fprintln(c.stdout, path)
		c.printTree(path, "", 0, *level, *dirOnly)
	}

	return nil
}

func (c *TreeCommand) printTree(root, prefix string, depth, maxDepth int, dirOnly bool) {
	// 检查深度限制
	if maxDepth > 0 && depth >= maxDepth {
		return
	}

	// 读取目录
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}

	// 过滤隐藏文件
	var filtered []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			if !dirOnly || entry.IsDir() {
				filtered = append(filtered, entry)
			}
		}
	}

	// 打印每个条目
	for i, entry := range filtered {
		isLast := i == len(filtered)-1

		// 确定分支字符
		var branch string
		if isLast {
			branch = "└── "
		} else {
			branch = "├── "
		}

		// 打印文件名
		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}
		fmt.Fprintf(c.stdout, "%s%s%s\n", prefix, branch, name)

		// 递归打印子目录
		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}

			subPath := filepath.Join(root, entry.Name())
			c.printTree(subPath, newPrefix, depth+1, maxDepth, dirOnly)
		}
	}
}

func (c *TreeCommand) Help() string {
	return `tree - 以树形结构显示目录

用法:
  tree [选项] [目录...]

选项:
  -L, --level        显示的最大层级深度
  -d, --directories  只显示目录

描述:
  以树形结构递归显示目录内容。

示例:
  tree                # 显示当前目录树
  tree -L 2           # 只显示 2 层
  tree -d             # 只显示目录
  tree src/           # 显示 src 目录树`
}

func (c *TreeCommand) ShortHelp() string {
	return "目录树显示"
}
