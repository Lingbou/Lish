package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type MvCommand struct {
	stdout *os.File
}

func NewMvCommand(stdout *os.File) *MvCommand {
	return &MvCommand{stdout: stdout}
}

func (c *MvCommand) Name() string {
	return "mv"
}

func (c *MvCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("mv", pflag.ContinueOnError)
	verbose := flags.BoolP("verbose", "v", false, "显示详细信息")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) < 2 {
		return fmt.Errorf("mv: 需要源文件和目标路径")
	}

	src := files[0]
	dst := files[1]

	// 检查源文件是否存在
	_, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("mv: %w", err)
	}

	// 检查目标是否是目录
	dstInfo, err := os.Stat(dst)
	if err == nil && dstInfo.IsDir() {
		// 目标是目录，保留原文件名
		dst = filepath.Join(dst, filepath.Base(src))
	}

	// 尝试直接重命名（同盘）
	err = os.Rename(src, dst)
	if err == nil {
		if *verbose {
			fmt.Fprintf(c.stdout, "'%s' -> '%s'\n", src, dst)
		}
		return nil
	}

	// 重命名失败，可能是跨盘，使用复制+删除
	if err := c.copyAndRemove(src, dst, *verbose); err != nil {
		return fmt.Errorf("mv: %w", err)
	}

	return nil
}

func (c *MvCommand) copyAndRemove(src, dst string, verbose bool) error {
	// 获取源信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return c.copyDirAndRemove(src, dst, verbose)
	}

	return c.copyFileAndRemove(src, dst, verbose)
}

func (c *MvCommand) copyFileAndRemove(src, dst string, verbose bool) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// 同步到磁盘
	if err := dstFile.Sync(); err != nil {
		return err
	}

	// 复制权限
	srcInfo, _ := srcFile.Stat()
	os.Chmod(dst, srcInfo.Mode())

	// 删除源文件
	if err := os.Remove(src); err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(c.stdout, "'%s' -> '%s'\n", src, dst)
	}

	return nil
}

func (c *MvCommand) copyDirAndRemove(src, dst string, verbose bool) error {
	// 获取源目录信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// 读取源目录内容
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 复制每个文件/目录
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := c.copyDirAndRemove(srcPath, dstPath, verbose); err != nil {
				return err
			}
		} else {
			if err := c.copyFileAndRemove(srcPath, dstPath, verbose); err != nil {
				return err
			}
		}
	}

	// 删除源目录
	if err := os.Remove(src); err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(c.stdout, "'%s' -> '%s'\n", src, dst)
	}

	return nil
}

func (c *MvCommand) Help() string {
	return `mv - 移动或重命名文件

用法:
  mv [选项] 源文件 目标文件
  mv [选项] 源文件... 目标目录

选项:
  -v, --verbose  显示详细信息

描述:
  将源文件移动到目标位置，或重命名文件。
  如果跨磁盘移动，会使用复制+删除的方式。

示例:
  mv file.txt newname.txt     # 重命名文件
  mv file.txt /path/to/dir/   # 移动到目录
  mv -v file1 file2 dir/      # 移动多个文件（详细模式）`
}

func (c *MvCommand) ShortHelp() string {
	return "移动或重命名文件"
}
