package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"

	flag "github.com/spf13/pflag"
)

// ChmodCommand chmod 命令 - 修改文件权限/属性
type ChmodCommand struct{}

func (c *ChmodCommand) Name() string {
	return "chmod"
}

func (c *ChmodCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("chmod", flag.ContinueOnError)
	recursive := flags.BoolP("recursive", "R", false, "递归修改目录")
	verbose := flags.BoolP("verbose", "v", false, "显示详细信息")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	if len(remaining) < 2 {
		return fmt.Errorf("用法: chmod [-R] [-v] <mode> <file>...")
	}

	mode := remaining[0]
	files := remaining[1:]

	// 解析模式
	perm, isWindows, err := parseMode(mode)
	if err != nil {
		return err
	}

	// 修改文件权限
	for _, file := range files {
		if err := chmodFile(file, perm, isWindows, *recursive, *verbose); err != nil {
			fmt.Fprintf(os.Stderr, "chmod: %s: %v\n", file, err)
		}
	}

	return nil
}

// parseMode 解析权限模式
func parseMode(mode string) (os.FileMode, bool, error) {
	isWindows := runtime.GOOS == "windows"

	// Windows 特殊模式
	if isWindows {
		switch mode {
		case "+r", "r":
			// 移除只读（让文件可读可写）
			return 0666, true, nil
		case "-r", "ro":
			// 设置只读
			return 0444, true, nil
		case "+w", "w", "rw":
			// 可读可写
			return 0666, true, nil
		case "-w":
			// 只读
			return 0444, true, nil
		}
	}

	// 尝试解析八进制权限（Unix 风格）
	if len(mode) == 3 || len(mode) == 4 {
		octal, err := strconv.ParseUint(mode, 8, 32)
		if err == nil {
			return os.FileMode(octal), false, nil
		}
	}

	// 符号模式（简化版）
	if len(mode) >= 2 {
		// u+x, a+w, o-r 等
		op := mode[1]
		perm := mode[2:]

		switch op {
		case '+':
			return parseSymbolicAdd(perm), false, nil
		case '-':
			return parseSymbolicRemove(perm), false, nil
		}
	}

	return 0, false, fmt.Errorf("无效的权限模式: %s", mode)
}

// parseSymbolicAdd 解析符号模式（添加权限）
func parseSymbolicAdd(perm string) os.FileMode {
	var mode os.FileMode = 0

	for _, p := range perm {
		switch p {
		case 'r':
			mode |= 0444 // 所有人可读
		case 'w':
			mode |= 0222 // 所有人可写
		case 'x':
			mode |= 0111 // 所有人可执行
		}
	}

	return mode
}

// parseSymbolicRemove 解析符号模式（移除权限）
func parseSymbolicRemove(perm string) os.FileMode {
	// 返回要移除的权限掩码
	var mode os.FileMode = 0777

	for _, p := range perm {
		switch p {
		case 'r':
			mode &^= 0444
		case 'w':
			mode &^= 0222
		case 'x':
			mode &^= 0111
		}
	}

	return mode
}

// chmodFile 修改文件权限
func chmodFile(path string, perm os.FileMode, isWindowsMode bool, recursive bool, verbose bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Windows 特殊处理
	if runtime.GOOS == "windows" && isWindowsMode {
		return chmodWindows(path, perm, verbose)
	}

	// Unix 风格权限
	if err := os.Chmod(path, perm); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("chmod: %s: 权限已修改为 %o\n", path, perm)
	}

	// 递归处理目录
	if recursive && info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			subPath := path + string(os.PathSeparator) + entry.Name()
			if err := chmodFile(subPath, perm, isWindowsMode, recursive, verbose); err != nil {
				fmt.Fprintf(os.Stderr, "chmod: %s: %v\n", subPath, err)
			}
		}
	}

	return nil
}

// chmodWindows Windows 文件属性处理
func chmodWindows(path string, perm os.FileMode, verbose bool) error {
	// Windows 上，我们使用 os.Chmod 来设置只读属性
	// perm == 0444 表示只读，0666 表示可读可写

	if err := os.Chmod(path, perm); err != nil {
		return err
	}

	if verbose {
		if perm == 0444 {
			fmt.Printf("chmod: %s: 已设置为只读\n", path)
		} else {
			fmt.Printf("chmod: %s: 已设置为可读写\n", path)
		}
	}

	return nil
}

func (c *ChmodCommand) Help() string {
	help := `chmod - 修改文件权限/属性

用法:
  chmod [-R] [-v] <mode> <file>...

说明:
  修改文件或目录的权限。
  `

	if runtime.GOOS == "windows" {
		help += `
Windows 模式:
  +r, r      移除只读（可读写）
  -r, ro     设置只读
  +w, w, rw  可读写
  -w         只读

  注意: Windows 不支持完整的 Unix 权限系统，
        chmod 主要用于设置只读属性。

示例:
  chmod +r file.txt      # 设置可读写
  chmod -r file.txt      # 设置只读
  chmod ro file.txt      # 设置只读
  chmod rw file.txt      # 设置可读写
  chmod -R +w dir/       # 递归设置可写
`
	} else {
		help += `
Unix 模式:
  数字模式:
    755      rwxr-xr-x
    644      rw-r--r--
    777      rwxrwxrwx
  
  符号模式:
    u+x      用户添加执行权限
    a+w      所有人添加写权限
    o-r      其他人移除读权限

选项:
  -R, --recursive    递归修改目录
  -v, --verbose      显示详细信息

示例:
  chmod 755 file         # 设置为 rwxr-xr-x
  chmod 644 file.txt     # 设置为 rw-r--r--
  chmod u+x script.sh    # 添加执行权限
  chmod -R 755 dir/      # 递归修改目录
`
	}

	return help
}

func (c *ChmodCommand) ShortHelp() string {
	if runtime.GOOS == "windows" {
		return "修改文件属性（只读/可写）"
	}
	return "修改文件权限"
}
