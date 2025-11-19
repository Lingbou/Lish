package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type RmCommand struct{}

func NewRmCommand() *RmCommand {
	return &RmCommand{}
}

func (c *RmCommand) Name() string {
	return "rm"
}

func (c *RmCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("rm", pflag.ContinueOnError)
	recursive := flags.BoolP("recursive", "r", false, "递归删除目录")
	force := flags.BoolP("force", "f", false, "强制删除，不提示")
	
	if err := flags.Parse(args); err != nil {
		return err
	}
	
	files := flags.Args()
	if len(files) == 0 {
		return fmt.Errorf("请指定至少一个文件或目录")
	}
	
	for _, file := range files {
		// 检查文件是否存在
		info, err := os.Stat(file)
		if err != nil {
			if *force && os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("访问 %s 失败: %w", file, err)
		}
		
		// 如果是目录但没有 -r 标志
		if info.IsDir() && !*recursive {
			return fmt.Errorf("%s 是目录，请使用 -r 选项", file)
		}
		
		// 删除文件或目录
		if *recursive {
			err = os.RemoveAll(file)
		} else {
			err = os.Remove(file)
		}
		
		if err != nil {
			return fmt.Errorf("删除 %s 失败: %w", file, err)
		}
	}
	
	return nil
}

func (c *RmCommand) Help() string {
	return `rm - 删除文件或目录

用法:
  rm [选项] 文件...

选项:
  -r, --recursive  递归删除目录及其内容
  -f, --force      强制删除，忽略不存在的文件

描述:
  删除指定的文件或目录。

警告:
  删除操作不可恢复，请谨慎使用。

示例:
  rm file.txt          # 删除文件
  rm -r dir            # 递归删除目录
  rm -rf dir           # 强制递归删除目录
  rm file1 file2       # 删除多个文件`
}

func (c *RmCommand) ShortHelp() string {
	return "删除文件或目录"
}

