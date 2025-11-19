package commands

import (
	"context"
	"fmt"
	"os"
	"time"
)

type TouchCommand struct{}

func NewTouchCommand() *TouchCommand {
	return &TouchCommand{}
}

func (c *TouchCommand) Name() string {
	return "touch"
}

func (c *TouchCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定至少一个文件名")
	}
	
	now := time.Now()
	
	for _, filename := range args {
		// 检查文件是否存在
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			// 创建新文件
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("创建文件 %s 失败: %w", filename, err)
			}
			file.Close()
		} else if err == nil {
			// 更新文件时间戳
			err = os.Chtimes(filename, now, now)
			if err != nil {
				return fmt.Errorf("更新文件 %s 时间戳失败: %w", filename, err)
			}
		} else {
			return fmt.Errorf("访问文件 %s 失败: %w", filename, err)
		}
	}
	
	return nil
}

func (c *TouchCommand) Help() string {
	return `touch - 创建空文件或更新文件时间戳

用法:
  touch 文件...

描述:
  如果文件不存在，创建一个空文件。
  如果文件存在，更新其访问和修改时间为当前时间。

示例:
  touch file.txt           # 创建或更新文件
  touch file1 file2 file3  # 创建或更新多个文件`
}

func (c *TouchCommand) ShortHelp() string {
	return "创建空文件或更新时间戳"
}

