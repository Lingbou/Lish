package commands

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

type GrepCommand struct {
	stdout *os.File
	stdin  *os.File
}

func NewGrepCommand(stdout, stdin *os.File) *GrepCommand {
	return &GrepCommand{
		stdout: stdout,
		stdin:  stdin,
	}
}

func (c *GrepCommand) Name() string {
	return "grep"
}

func (c *GrepCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("grep", pflag.ContinueOnError)
	ignoreCase := flags.BoolP("ignore-case", "i", false, "忽略大小写")
	lineNumber := flags.BoolP("line-number", "n", false, "显示行号")
	invert := flags.BoolP("invert-match", "v", false, "反向匹配")
	recursive := flags.BoolP("recursive", "r", false, "递归搜索目录")

	if err := flags.Parse(args); err != nil {
		return err
	}

	grepArgs := flags.Args()
	if len(grepArgs) == 0 {
		return fmt.Errorf("grep: 需要指定搜索模式")
	}

	pattern := grepArgs[0]
	files := grepArgs[1:]

	// 编译正则表达式
	regexFlags := ""
	if *ignoreCase {
		regexFlags = "(?i)"
	}
	re, err := regexp.Compile(regexFlags + pattern)
	if err != nil {
		return fmt.Errorf("grep: 无效的正则表达式: %w", err)
	}

	// 如果没有指定文件，从标准输入读取
	if len(files) == 0 {
		return c.grepReader(c.stdin, "-", re, *lineNumber, *invert)
	}

	// 搜索指定文件
	for _, file := range files {
		if *recursive {
			if err := c.grepRecursive(file, re, *lineNumber, *invert); err != nil {
				fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			}
		} else {
			if err := c.grepFile(file, re, *lineNumber, *invert, len(files) > 1); err != nil {
				fmt.Fprintf(os.Stderr, "grep: %v\n", err)
			}
		}
	}

	return nil
}

func (c *GrepCommand) grepFile(filename string, re *regexp.Regexp, showLine, invert, showFilename bool) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := re.MatchString(line)

		// 反向匹配
		if invert {
			matched = !matched
		}

		if matched {
			c.printMatch(filename, lineNum, line, showLine, showFilename, re)
		}
	}

	return scanner.Err()
}

func (c *GrepCommand) grepReader(reader io.Reader, name string, re *regexp.Regexp, showLine, invert bool) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		matched := re.MatchString(line)

		if invert {
			matched = !matched
		}

		if matched {
			c.printMatch(name, lineNum, line, showLine, false, re)
		}
	}

	return scanner.Err()
}

func (c *GrepCommand) grepRecursive(path string, re *regexp.Regexp, showLine, invert bool) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略无法访问的文件
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 跳过二进制文件（简单检测）
		if !c.isTextFile(path) {
			return nil
		}

		return c.grepFile(path, re, showLine, invert, true)
	})
}

func (c *GrepCommand) isTextFile(path string) bool {
	// 简单的文本文件检测
	ext := strings.ToLower(filepath.Ext(path))
	textExts := []string{".txt", ".md", ".go", ".js", ".py", ".java", ".c", ".cpp", ".h", ".html", ".css", ".json", ".xml", ".yaml", ".yml"}

	for _, te := range textExts {
		if ext == te {
			return true
		}
	}

	// 如果没有扩展名，尝试读取前几个字节判断
	if ext == "" {
		file, err := os.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()

		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		return c.isPrintable(buf[:n])
	}

	return false
}

func (c *GrepCommand) isPrintable(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return true
}

func (c *GrepCommand) printMatch(filename string, lineNum int, line string, showLine, showFilename bool, re *regexp.Regexp) {
	var prefix string

	if showFilename {
		prefix = fmt.Sprintf("%s:", filename)
	}

	if showLine {
		prefix += fmt.Sprintf("%d:", lineNum)
	}

	// 彩色高亮匹配部分
	const (
		colorRed   = "\033[31m"
		colorReset = "\033[0m"
	)

	highlighted := re.ReplaceAllStringFunc(line, func(match string) string {
		return colorRed + match + colorReset
	})

	fmt.Fprintf(c.stdout, "%s%s\n", prefix, highlighted)
}

func (c *GrepCommand) Help() string {
	return `grep - 搜索文本

用法:
  grep [选项] 模式 [文件...]

选项:
  -i, --ignore-case   忽略大小写
  -n, --line-number   显示行号
  -v, --invert-match  反向匹配（显示不匹配的行）
  -r, --recursive     递归搜索目录

描述:
  在文件中搜索匹配模式的行。支持正则表达式。
  如果不指定文件，从标准输入读取。

示例:
  grep "error" log.txt            # 搜索 "error"
  grep -i "ERROR" log.txt         # 忽略大小写搜索
  grep -n "error" log.txt         # 显示行号
  grep -r "TODO" .                # 递归搜索当前目录
  ls | grep ".txt"                # 从管道输入搜索
  grep -v "^#" file.txt           # 显示非注释行`
}

func (c *GrepCommand) ShortHelp() string {
	return "搜索文本"
}
