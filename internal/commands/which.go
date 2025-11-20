package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type WhichCommand struct {
	stdout   *os.File
	registry *Registry
}

func NewWhichCommand(stdout *os.File, registry *Registry) *WhichCommand {
	return &WhichCommand{
		stdout:   stdout,
		registry: registry,
	}
}

func (c *WhichCommand) Name() string {
	return "which"
}

func (c *WhichCommand) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("which: 需要指定命令名")
	}
	
	for _, cmdName := range args {
		// 先检查是否是内置命令
		if _, exists := c.registry.Get(cmdName); exists {
			fmt.Fprintf(c.stdout, "%s: Lish 内置命令\n", cmdName)
			continue
		}
		
		// 在 PATH 中查找
		path, err := c.findInPath(cmdName)
		if err != nil {
			fmt.Fprintf(c.stdout, "%s: 未找到\n", cmdName)
		} else {
			fmt.Fprintln(c.stdout, path)
		}
	}
	
	return nil
}

func (c *WhichCommand) findInPath(cmdName string) (string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH 环境变量未设置")
	}
	
	// Windows 下需要添加 .exe 扩展名
	extensions := []string{""}
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		extensions = []string{".exe", ".bat", ".cmd", ""}
	}
	
	paths := filepath.SplitList(pathEnv)
	for _, dir := range paths {
		for _, ext := range extensions {
			fullPath := filepath.Join(dir, cmdName+ext)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}
	
	return "", fmt.Errorf("未找到命令")
}

func (c *WhichCommand) Help() string {
	return `which - 查找命令路径

用法:
  which 命令名...

描述:
  在 PATH 环境变量中查找命令的完整路径。
  对于 Lish 内置命令，会显示 "内置命令"。

示例:
  which ls               # 查找 ls 命令
  which cmd              # 查找 cmd 命令
  which git code         # 查找多个命令`
}

func (c *WhichCommand) ShortHelp() string {
	return "查找命令路径"
}

