package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type FindCommand struct {
	stdout *os.File
}

func NewFindCommand(stdout *os.File) *FindCommand {
	return &FindCommand{stdout: stdout}
}

func (c *FindCommand) Name() string {
	return "find"
}

func (c *FindCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("find", pflag.ContinueOnError)
	namePattern := flags.StringP("name", "n", "", "按文件名查找（支持通配符）")
	fileType := flags.StringP("type", "t", "", "按类型查找（f=文件, d=目录）")

	if err := flags.Parse(args); err != nil {
		return err
	}

	paths := flags.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		if err := c.findInPath(path, *namePattern, *fileType); err != nil {
			fmt.Fprintf(os.Stderr, "find: %v\n", err)
		}
	}

	return nil
}

func (c *FindCommand) findInPath(root, namePattern, fileType string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略无法访问的文件
		}

		// 过滤文件类型
		if fileType != "" {
			if fileType == "f" && info.IsDir() {
				return nil
			}
			if fileType == "d" && !info.IsDir() {
				return nil
			}
		}

		// 过滤文件名
		if namePattern != "" {
			matched, err := filepath.Match(namePattern, filepath.Base(path))
			if err != nil || !matched {
				return nil
			}
		}

		// 打印匹配的路径
		fmt.Fprintln(c.stdout, path)

		return nil
	})
}

func (c *FindCommand) Help() string {
	return `find - 查找文件

用法:
  find [路径...] [选项]

选项:
  -n, --name  按文件名查找（支持通配符 * ?）
  -t, --type  按类型查找（f=文件, d=目录）

描述:
  在指定目录中查找文件。如果不指定路径，默认在当前目录查找。

示例:
  find                    # 列出当前目录所有文件
  find /path              # 列出指定目录所有文件
  find -n "*.txt"         # 查找所有 .txt 文件
  find -t f               # 只查找文件
  find -t d               # 只查找目录
  find -n "test*" -t f    # 查找以 test 开头的文件`
}

func (c *FindCommand) ShortHelp() string {
	return "查找文件"
}
