package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type CdCommand struct {
	lastDir string
}

func NewCdCommand() *CdCommand {
	return &CdCommand{}
}

func (c *CdCommand) Name() string {
	return "cd"
}

func (c *CdCommand) Execute(ctx context.Context, args []string) error {
	var targetDir string

	if len(args) == 0 {
		// 无参数时，切换到用户目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("获取用户目录失败: %w", err)
		}
		targetDir = homeDir
	} else {
		arg := args[0]

		switch arg {
		case "-":
			// 切换到上次的目录
			if c.lastDir == "" {
				return fmt.Errorf("没有上次访问的目录")
			}
			targetDir = c.lastDir
		case "~":
			// 切换到用户目录
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("获取用户目录失败: %w", err)
			}
			targetDir = homeDir
		default:
			// 处理 ~ 开头的路径
			if len(arg) > 0 && arg[0] == '~' {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("获取用户目录失败: %w", err)
				}
				targetDir = filepath.Join(homeDir, arg[1:])
			} else {
				targetDir = arg
			}
		}
	}

	// 获取当前目录（用于保存）
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}

	// 切换目录
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("切换目录失败: %w", err)
	}

	// 保存上次的目录
	c.lastDir = currentDir

	return nil
}

func (c *CdCommand) Help() string {
	return `cd - 切换目录

用法:
  cd [目录]

描述:
  切换当前工作目录。

参数:
  目录    目标目录路径
          无参数时切换到用户主目录
          ~ 或 ~/ 表示用户主目录
          - 表示上次访问的目录

示例:
  cd               # 切换到用户主目录
  cd /root         # 切换到 /root
  cd ~             # 切换到用户主目录
  cd ~/Documents   # 切换到用户主目录下的 Documents
  cd -             # 切换到上次访问的目录
  cd ..            # 切换到上级目录`
}

func (c *CdCommand) ShortHelp() string {
	return "切换目录"
}
