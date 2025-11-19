package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type CpCommand struct {
	stdout *os.File
}

func NewCpCommand(stdout *os.File) *CpCommand {
	return &CpCommand{stdout: stdout}
}

func (c *CpCommand) Name() string {
	return "cp"
}

func (c *CpCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("cp", pflag.ContinueOnError)
	recursive := flags.BoolP("recursive", "r", false, "递归复制目录")
	verbose := flags.BoolP("verbose", "v", false, "显示详细信息")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) < 2 {
		return fmt.Errorf("cp: 需要源文件和目标路径")
	}

	src := files[0]
	dst := files[1]

	// 检查源文件
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cp: %w", err)
	}

	// 如果是目录但没有 -r 标志
	if srcInfo.IsDir() && !*recursive {
		return fmt.Errorf("cp: %s 是目录（省略目录，使用 -r 递归复制）", src)
	}

	// 复制
	if srcInfo.IsDir() {
		return c.copyDir(src, dst, *verbose)
	}

	return c.copyFile(src, dst, *verbose)
}

func (c *CpCommand) copyFile(src, dst string, verbose bool) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	// 检查目标是否是目录
	dstInfo, err := os.Stat(dst)
	if err == nil && dstInfo.IsDir() {
		// 目标是目录，保留原文件名
		dst = filepath.Join(dst, filepath.Base(src))
	}

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 复制权限
	srcInfo, _ := srcFile.Stat()
	os.Chmod(dst, srcInfo.Mode())

	if verbose {
		fmt.Fprintf(c.stdout, "'%s' -> '%s'\n", src, dst)
	}

	return nil
}

func (c *CpCommand) copyDir(src, dst string, verbose bool) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	if verbose {
		fmt.Fprintf(c.stdout, "'%s' -> '%s'\n", src, dst)
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 递归复制每个文件/目录
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := c.copyDir(srcPath, dstPath, verbose); err != nil {
				return err
			}
		} else {
			if err := c.copyFile(srcPath, dstPath, verbose); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CpCommand) Help() string {
	return `cp - 复制文件或目录

用法:
  cp [选项] 源文件 目标文件
  cp [选项] 源文件... 目标目录

选项:
  -r, --recursive  递归复制目录及其内容
  -v, --verbose    显示详细信息

描述:
  将源文件或目录复制到目标位置。

示例:
  cp file.txt backup.txt        # 复制文件
  cp file.txt /path/to/dir/     # 复制到目录
  cp -r dir1 dir2               # 递归复制目录
  cp -v file1 file2 file3 dir/  # 复制多个文件（详细模式）`
}

func (c *CpCommand) ShortHelp() string {
	return "复制文件或目录"
}
