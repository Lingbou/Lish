package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// ChownCommand chown 命令 - 修改文件所有者
type ChownCommand struct{}

func (c *ChownCommand) Name() string {
	return "chown"
}

func (c *ChownCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("chown", flag.ContinueOnError)
	recursive := flags.BoolP("recursive", "R", false, "递归修改目录")
	verbose := flags.BoolP("verbose", "v", false, "显示详细信息")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	if len(remaining) < 2 {
		return fmt.Errorf("用法: chown [-R] [-v] [user][:group] <file>...")
	}

	owner := remaining[0]
	files := remaining[1:]

	// Windows 平台提示
	if runtime.GOOS == "windows" {
		fmt.Println("⚠️  注意: Windows 平台的 chown 功能有限")
		fmt.Println("   修改文件所有者需要管理员权限")
		fmt.Println("   当前版本仅显示文件信息，不执行实际修改")
		fmt.Println()
	}

	// 解析所有者信息
	user, group := parseOwner(owner)

	// 修改文件所有者
	for _, file := range files {
		if err := chownFile(file, user, group, *recursive, *verbose); err != nil {
			fmt.Fprintf(os.Stderr, "chown: %s: %v\n", file, err)
		}
	}

	return nil
}

// parseOwner 解析所有者信息
func parseOwner(owner string) (string, string) {
	parts := strings.Split(owner, ":")
	user := parts[0]
	group := ""
	if len(parts) > 1 {
		group = parts[1]
	}
	return user, group
}

// chownFile 修改文件所有者
func chownFile(path string, user string, group string, recursive bool, verbose bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Windows 平台处理
	if runtime.GOOS == "windows" {
		return chownWindows(path, user, group)
	}

	// Unix 平台处理
	return chownUnix(path, user, group, info, recursive, verbose)
}

// chownWindows Windows 平台处理
func chownWindows(path string, user string, group string) error {
	// Windows 上修改文件所有者需要使用 Windows API
	// 这里只显示信息，不执行实际修改
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	fmt.Printf("%s:\n", path)
	fmt.Printf("  类型: ")
	if info.IsDir() {
		fmt.Println("目录")
	} else {
		fmt.Println("文件")
	}
	fmt.Printf("  大小: %d 字节\n", info.Size())
	fmt.Printf("  修改时间: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))
	fmt.Printf("  请求更改所有者为: %s", user)
	if group != "" {
		fmt.Printf(":%s", group)
	}
	fmt.Println()
	fmt.Println("  (Windows 平台需要管理员权限，当前未执行)")
	fmt.Println()

	return nil
}

// chownUnix Unix 平台处理
func chownUnix(path string, user string, group string, info os.FileInfo, recursive bool, verbose bool) error {
	// 解析 UID 和 GID
	uid := -1
	gid := -1

	if user != "" {
		// 尝试解析为数字 UID
		if id, err := strconv.Atoi(user); err == nil {
			uid = id
		} else {
			// TODO: 通过用户名查找 UID（需要额外的系统调用）
			return fmt.Errorf("暂不支持用户名解析，请使用数字 UID")
		}
	}

	if group != "" {
		// 尝试解析为数字 GID
		if id, err := strconv.Atoi(group); err == nil {
			gid = id
		} else {
			return fmt.Errorf("暂不支持组名解析，请使用数字 GID")
		}
	}

	// 执行 chown
	if err := os.Chown(path, uid, gid); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("chown: %s: 所有者已修改为 %s", path, user)
		if group != "" {
			fmt.Printf(":%s", group)
		}
		fmt.Println()
	}

	// 递归处理目录
	if recursive && info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			subPath := path + string(os.PathSeparator) + entry.Name()
			subInfo, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "chown: %s: %v\n", subPath, err)
				continue
			}
			if err := chownUnix(subPath, user, group, subInfo, recursive, verbose); err != nil {
				fmt.Fprintf(os.Stderr, "chown: %s: %v\n", subPath, err)
			}
		}
	}

	return nil
}

func (c *ChownCommand) Help() string {
	help := `chown - 修改文件所有者

用法:
  chown [-R] [-v] [user][:group] <file>...

说明:
  修改文件或目录的所有者和/或所属组。
`

	if runtime.GOOS == "windows" {
		help += `
Windows 平台说明:
  • 修改文件所有者需要管理员权限
  • 当前版本仅显示文件信息
  • 不执行实际的所有者修改操作
  • 建议使用 Windows 文件属性界面进行修改

选项:
  -R, --recursive    递归修改目录
  -v, --verbose      显示详细信息

示例:
  chown user file.txt           # 显示文件信息
  chown user:group file.txt     # 显示文件信息
  chown -R user dir/            # 递归显示目录信息
`
	} else {
		help += `
选项:
  -R, --recursive    递归修改目录
  -v, --verbose      显示详细信息

示例:
  chown 1000 file.txt           # 修改所有者 UID
  chown 1000:1000 file.txt      # 修改所有者和组
  chown -R 1000 dir/            # 递归修改目录

注意:
  • 需要 root 权限或文件所有者权限
  • 当前版本使用数字 UID/GID
  • 不支持用户名/组名解析
`
	}

	return help
}

func (c *ChownCommand) ShortHelp() string {
	if runtime.GOOS == "windows" {
		return "查看文件所有者信息"
	}
	return "修改文件所有者"
}

