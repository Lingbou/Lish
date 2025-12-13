package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"

	flag "github.com/spf13/pflag"
)

// LnCommand ln 命令 - 创建链接
type LnCommand struct{}

func (c *LnCommand) Name() string {
	return "ln"
}

func (c *LnCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("ln", flag.ContinueOnError)
	symbolic := flags.BoolP("symbolic", "s", false, "创建符号链接")
	force := flags.BoolP("force", "f", false, "强制覆盖已存在的链接")
	verbose := flags.BoolP("verbose", "v", false, "显示详细信息")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	if len(remaining) < 2 {
		return fmt.Errorf("用法: ln [-s] [-f] [-v] <target> <link_name>")
	}

	target := remaining[0]
	linkName := remaining[1]

	// Windows 平台提示
	if runtime.GOOS == "windows" {
		fmt.Println("⚠️  注意: Windows 平台创建符号链接的要求")
		fmt.Println("   • 需要管理员权限，或")
		fmt.Println("   • 启用开发者模式")
		fmt.Println("   • 硬链接只支持文件，不支持目录")
		fmt.Println()
	}

	// 检查目标是否存在
	if _, err := os.Stat(target); os.IsNotExist(err) {
		fmt.Printf("⚠️  警告: 目标 '%s' 不存在\n", target)
		if *symbolic {
			fmt.Println("   符号链接允许指向不存在的目标")
		} else {
			return fmt.Errorf("硬链接的目标必须存在")
		}
	}

	// 检查链接是否已存在
	if _, err := os.Lstat(linkName); err == nil {
		if *force {
			if err := os.Remove(linkName); err != nil {
				return fmt.Errorf("无法删除已存在的链接: %w", err)
			}
			if *verbose {
				fmt.Printf("已删除已存在的链接: %s\n", linkName)
			}
		} else {
			return fmt.Errorf("链接已存在: %s (使用 -f 强制覆盖)", linkName)
		}
	}

	// 创建链接
	var err error
	if *symbolic {
		// 创建符号链接
		err = os.Symlink(target, linkName)
		if err != nil {
			if runtime.GOOS == "windows" {
				return fmt.Errorf("创建符号链接失败: %w\n提示: 请以管理员身份运行，或启用开发者模式", err)
			}
			return fmt.Errorf("创建符号链接失败: %w", err)
		}
		if *verbose {
			fmt.Printf("✓ 已创建符号链接: %s -> %s\n", linkName, target)
		}
	} else {
		// 创建硬链接
		err = os.Link(target, linkName)
		if err != nil {
			return fmt.Errorf("创建硬链接失败: %w", err)
		}
		if *verbose {
			fmt.Printf("✓ 已创建硬链接: %s -> %s\n", linkName, target)
		}
	}

	return nil
}

func (c *LnCommand) Help() string {
	help := `ln - 创建文件链接

用法:
  ln [-s] [-f] [-v] <target> <link_name>

说明:
  创建文件或目录的链接。支持硬链接和符号链接。
`

	if runtime.GOOS == "windows" {
		help += `
Windows 平台说明:
  • 符号链接需要管理员权限或开发者模式
  • 硬链接只支持文件，不支持目录
  • 符号链接支持相对路径和绝对路径
  
  启用开发者模式:
    设置 -> 更新和安全 -> 开发者选项 -> 开发者模式

选项:
  -s, --symbolic     创建符号链接（推荐）
  -f, --force        强制覆盖已存在的链接
  -v, --verbose      显示详细信息

示例:
  ln -s target.txt link.txt         # 创建符号链接
  ln -s C:\data\file.txt link.txt   # 绝对路径符号链接
  ln -s dir linkdir                 # 目录符号链接
  ln target.txt hardlink.txt        # 硬链接（文件）
  ln -sf target.txt link.txt        # 强制覆盖

链接类型区别:
  符号链接 (-s):
    • 指向路径（可以跨分区）
    • 删除原文件后链接失效
    • 支持目录
    • 可以查看链接目标
  
  硬链接:
    • 指向文件内容（不能跨分区）
    • 删除原文件后仍然有效
    • 只支持文件
    • 多个硬链接共享相同数据
`
	} else {
		help += `
选项:
  -s, --symbolic     创建符号链接
  -f, --force        强制覆盖已存在的链接
  -v, --verbose      显示详细信息

示例:
  ln -s target.txt link.txt         # 创建符号链接
  ln -s /path/to/file link          # 绝对路径符号链接
  ln -s ../data/file.txt link.txt   # 相对路径符号链接
  ln target.txt hardlink.txt        # 创建硬链接
  ln -sf target.txt link.txt        # 强制覆盖已存在的链接

链接类型:
  符号链接 (-s):
    • 类似于快捷方式
    • 可以指向文件或目录
    • 可以跨文件系统
    • 删除原文件后链接失效
  
  硬链接:
    • 同一文件的另一个名称
    • 只能指向文件
    • 不能跨文件系统
    • 删除原文件不影响硬链接
`
	}

	return help
}

func (c *LnCommand) ShortHelp() string {
	return "创建文件链接"
}
