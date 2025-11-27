package commands

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type UnzipCommand struct {
	stdout *os.File
}

func NewUnzipCommand(stdout *os.File) *UnzipCommand {
	return &UnzipCommand{stdout: stdout}
}

func (c *UnzipCommand) Name() string {
	return "unzip"
}

func (c *UnzipCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("unzip", pflag.ContinueOnError)
	outputDir := flags.StringP("dir", "d", ".", "解压到指定目录")
	list := flags.BoolP("list", "l", false, "列出压缩包内容")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) == 0 {
		return fmt.Errorf("unzip: 需要指定 ZIP 文件")
	}

	zipFile := files[0]

	// 打开 zip 文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("打开压缩包失败: %w", err)
	}
	defer reader.Close()

	// 如果只是列出内容
	if *list {
		return c.listZip(reader)
	}

	// 解压文件
	for _, file := range reader.File {
		if err := c.extractFile(file, *outputDir); err != nil {
			return err
		}
	}

	fmt.Fprintf(c.stdout, "✓ 已解压到: %s\n", *outputDir)
	return nil
}

func (c *UnzipCommand) listZip(reader *zip.ReadCloser) error {
	fmt.Fprintln(c.stdout, "压缩包内容:")
	for _, file := range reader.File {
		fmt.Fprintf(c.stdout, "  %s (%d bytes)\n", file.Name, file.UncompressedSize64)
	}
	return nil
}

func (c *UnzipCommand) extractFile(file *zip.File, outputDir string) error {
	// 构建输出路径
	outPath := filepath.Join(outputDir, file.Name)

	// 如果是目录
	if file.FileInfo().IsDir() {
		return os.MkdirAll(outPath, file.Mode())
	}

	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	// 打开压缩文件
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制内容
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (c *UnzipCommand) Help() string {
	return `unzip - 解压 ZIP 文件

用法:
  unzip [选项] 压缩包.zip

选项:
  -d, --dir   解压到指定目录（默认当前目录）
  -l, --list  列出压缩包内容（不解压）

描述:
  解压 ZIP 格式的压缩包。

示例:
  unzip archive.zip                  # 解压到当前目录
  unzip archive.zip -d /tmp          # 解压到 /tmp
  unzip -l archive.zip               # 列出内容`
}

func (c *UnzipCommand) ShortHelp() string {
	return "解压 ZIP 文件"
}
