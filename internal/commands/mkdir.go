package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type MkdirCommand struct{}

func NewMkdirCommand() *MkdirCommand {
	return &MkdirCommand{}
}

func (c *MkdirCommand) Name() string {
	return "mkdir"
}

func (c *MkdirCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("mkdir", pflag.ContinueOnError)
	parents := flags.BoolP("parents", "p", false, "递归创建父目录")
	
	if err := flags.Parse(args); err != nil {
		return err
	}
	
	dirs := flags.Args()
	if len(dirs) == 0 {
		return fmt.Errorf("请指定至少一个目录名")
	}
	
	for _, dir := range dirs {
		var err error
		if *parents {
			err = os.MkdirAll(dir, 0755)
		} else {
			err = os.Mkdir(dir, 0755)
		}
		
		if err != nil {
			return fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}
	
	return nil
}

func (c *MkdirCommand) Help() string {
	return `mkdir - 创建目录

用法:
  mkdir [选项] 目录...

选项:
  -p, --parents  递归创建目录，如果父目录不存在则创建

描述:
  创建一个或多个目录。

示例:
  mkdir newdir              # 创建单个目录
  mkdir dir1 dir2 dir3      # 创建多个目录
  mkdir -p parent/child     # 递归创建目录`
}

func (c *MkdirCommand) ShortHelp() string {
	return "创建目录"
}

