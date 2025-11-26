package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type DiffCommand struct {
	stdout *os.File
}

func NewDiffCommand(stdout *os.File) *DiffCommand {
	return &DiffCommand{stdout: stdout}
}

func (c *DiffCommand) Name() string {
	return "diff"
}

func (c *DiffCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("diff", pflag.ContinueOnError)
	brief := flags.BoolP("brief", "q", false, "只显示文件是否不同")

	if err := flags.Parse(args); err != nil {
		return err
	}

	files := flags.Args()
	if len(files) < 2 {
		return fmt.Errorf("diff: 需要两个文件")
	}

	file1, file2 := files[0], files[1]

	// 读取两个文件
	lines1, err := c.readLines(file1)
	if err != nil {
		return fmt.Errorf("读取 %s 失败: %w", file1, err)
	}

	lines2, err := c.readLines(file2)
	if err != nil {
		return fmt.Errorf("读取 %s 失败: %w", file2, err)
	}

	// 比较文件
	hasDiff := c.compareFiles(lines1, lines2, file1, file2, *brief)

	if !hasDiff {
		if !*brief {
			fmt.Fprintln(c.stdout, "文件相同")
		}
	}

	return nil
}

func (c *DiffCommand) readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func (c *DiffCommand) compareFiles(lines1, lines2 []string, file1, file2 string, brief bool) bool {
	maxLen := len(lines1)
	if len(lines2) > maxLen {
		maxLen = len(lines2)
	}

	hasDiff := false

	for i := 0; i < maxLen; i++ {
		line1 := ""
		if i < len(lines1) {
			line1 = lines1[i]
		}

		line2 := ""
		if i < len(lines2) {
			line2 = lines2[i]
		}

		if line1 != line2 {
			hasDiff = true

			if brief {
				fmt.Fprintf(c.stdout, "文件 %s 和 %s 不同\n", file1, file2)
				return true
			}

			// 显示差异
			lineNum := i + 1
			if i < len(lines1) {
				fmt.Fprintf(c.stdout, "< %d: %s\n", lineNum, line1)
			}
			if i < len(lines2) {
				fmt.Fprintf(c.stdout, "> %d: %s\n", lineNum, line2)
			}
			fmt.Fprintln(c.stdout, "---")
		}
	}

	return hasDiff
}

func (c *DiffCommand) Help() string {
	return `diff - 比较文件差异

用法:
  diff [选项] 文件1 文件2

选项:
  -q, --brief  只显示文件是否不同

描述:
  逐行比较两个文件的差异。

示例:
  diff file1.txt file2.txt     # 显示详细差异
  diff -q file1.txt file2.txt  # 只显示是否不同`
}

func (c *DiffCommand) ShortHelp() string {
	return "比较文件差异"
}
