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

type ZipCommand struct {
	stdout *os.File
}

func NewZipCommand(stdout *os.File) *ZipCommand {
	return &ZipCommand{stdout: stdout}
}

func (c *ZipCommand) Name() string {
	return "zip"
}

func (c *ZipCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("zip", pflag.ContinueOnError)
	recursive := flags.BoolP("recursive", "r", false, "递归压缩目录")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) < 2 {
		return fmt.Errorf("zip: 用法: zip [-r] 压缩包名.zip 文件...")
	}

	zipFile := files[0]
	sources := files[1:]

	// 创建 zip 文件
	archive, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("创建压缩包失败: %w", err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// 添加文件
	for _, source := range sources {
		if err := c.addToZip(zipWriter, source, *recursive); err != nil {
			return err
		}
	}

	fmt.Fprintf(c.stdout, "✓ 已创建压缩包: %s\n", zipFile)
	return nil
}

func (c *ZipCommand) addToZip(zipWriter *zip.Writer, source string, recursive bool) error {
	info, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("访问 %s 失败: %w", source, err)
	}

	if info.IsDir() {
		if !recursive {
			return fmt.Errorf("跳过目录 %s（使用 -r 递归压缩）", source)
		}
		return c.addDirToZip(zipWriter, source)
	}

	return c.addFileToZip(zipWriter, source)
}

func (c *ZipCommand) addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func (c *ZipCommand) addDirToZip(zipWriter *zip.Writer, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		writer, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}

func (c *ZipCommand) Help() string {
	return `zip - 创建 ZIP 压缩包

用法:
  zip [选项] 压缩包名.zip 文件或目录...

选项:
  -r, --recursive  递归压缩目录

描述:
  将文件或目录压缩为 ZIP 格式。

示例:
  zip archive.zip file.txt              # 压缩单个文件
  zip archive.zip file1.txt file2.txt   # 压缩多个文件
  zip -r backup.zip mydir/              # 递归压缩目录`
}

func (c *ZipCommand) ShortHelp() string {
	return "创建 ZIP 压缩包"
}
