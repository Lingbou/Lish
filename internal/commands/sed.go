package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

// SedCommand sed 命令 - 流编辑器（基础版）
type SedCommand struct{}

func (c *SedCommand) Name() string {
	return "sed"
}

func (c *SedCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("sed", flag.ContinueOnError)
	inPlace := flags.BoolP("in-place", "i", false, "原地编辑文件")
	quiet := flags.BoolP("quiet", "n", false, "静默模式（不自动打印）")
	expression := flags.StringP("expression", "e", "", "指定sed表达式")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	// 确定表达式和文件
	var script string
	var filenames []string

	if *expression != "" {
		script = *expression
		filenames = remaining
	} else if len(remaining) > 0 {
		script = remaining[0]
		filenames = remaining[1:]
	} else {
		return fmt.Errorf("用法: sed [-i] [-n] [-e expression] 'script' [file...]")
	}

	// 读取输入
	var lines []string
	var err error

	if len(filenames) == 0 {
		// 从标准输入读取
		lines, err = readLines(os.Stdin)
		if err != nil {
			return fmt.Errorf("读取输入失败: %w", err)
		}
	} else {
		// 从文件读取
		for _, filename := range filenames {
			fileLines, err := readLinesFromFile(filename)
			if err != nil {
				return fmt.Errorf("读取文件 %s 失败: %w", filename, err)
			}
			lines = append(lines, fileLines...)
		}
	}

	// 执行sed脚本
	result, err := executeSedScript(script, lines, *quiet)
	if err != nil {
		return err
	}

	// 输出或写回文件
	if *inPlace && len(filenames) > 0 {
		// 原地编辑模式
		return writeLinesToFile(filenames[0], result)
	} else {
		// 输出到标准输出
		for _, line := range result {
			fmt.Println(line)
		}
	}

	return nil
}

// executeSedScript 执行sed脚本
func executeSedScript(script string, lines []string, quiet bool) ([]string, error) {
	var result []string

	// 解析脚本命令
	cmd, err := parseSedCommand(script)
	if err != nil {
		return nil, err
	}

	for lineNum, line := range lines {
		processedLine, shouldPrint, err := processSedCommand(cmd, line, lineNum)
		if err != nil {
			return nil, err
		}

		// 根据命令决定是否输出
		if shouldPrint || (!quiet && cmd.action != "d") {
			result = append(result, processedLine)
		}
	}

	return result, nil
}

// sedCommand sed命令结构
type sedCommand struct {
	action      string      // s, d, p, a, i
	pattern     string      // 匹配模式
	replacement string      // 替换字符串
	flags       string      // g, i等标志
	regex       *regexp.Regexp
}

// parseSedCommand 解析sed命令
func parseSedCommand(script string) (*sedCommand, error) {
	cmd := &sedCommand{}

	script = strings.TrimSpace(script)

	// 替换命令: s/pattern/replacement/flags
	if strings.HasPrefix(script, "s/") || strings.HasPrefix(script, "s|") {
		delimiter := script[1]
		parts := strings.Split(script[2:], string(delimiter))
		if len(parts) < 2 {
			return nil, fmt.Errorf("无效的替换命令: %s", script)
		}

		cmd.action = "s"
		cmd.pattern = parts[0]
		cmd.replacement = parts[1]
		if len(parts) > 2 {
			cmd.flags = parts[2]
		}

		// 编译正则表达式
		regexFlags := ""
		if strings.Contains(cmd.flags, "i") {
			regexFlags = "(?i)"
		}
		regex, err := regexp.Compile(regexFlags + cmd.pattern)
		if err != nil {
			return nil, fmt.Errorf("无效的正则表达式: %w", err)
		}
		cmd.regex = regex

		return cmd, nil
	}

	// 删除命令: d
	if script == "d" {
		cmd.action = "d"
		return cmd, nil
	}

	// 打印命令: p
	if script == "p" {
		cmd.action = "p"
		return cmd, nil
	}

	// 追加命令: a\text
	if strings.HasPrefix(script, "a\\") {
		cmd.action = "a"
		cmd.replacement = strings.TrimPrefix(script, "a\\")
		return cmd, nil
	}

	// 插入命令: i\text
	if strings.HasPrefix(script, "i\\") {
		cmd.action = "i"
		cmd.replacement = strings.TrimPrefix(script, "i\\")
		return cmd, nil
	}

	return nil, fmt.Errorf("不支持的sed命令: %s", script)
}

// processSedCommand 处理sed命令
func processSedCommand(cmd *sedCommand, line string, lineNum int) (string, bool, error) {
	switch cmd.action {
	case "s":
		// 替换
		if strings.Contains(cmd.flags, "g") {
			// 全局替换
			return cmd.regex.ReplaceAllString(line, cmd.replacement), false, nil
		}
		// 只替换第一个
		return cmd.regex.ReplaceAllStringFunc(line, func(match string) string {
			return cmd.replacement
		}), false, nil

	case "d":
		// 删除（不输出）
		return "", false, nil

	case "p":
		// 打印（额外输出一次）
		return line, true, nil

	case "a":
		// 追加（在行后添加）
		return line + "\n" + cmd.replacement, false, nil

	case "i":
		// 插入（在行前添加）
		return cmd.replacement + "\n" + line, false, nil

	default:
		return line, false, fmt.Errorf("未知命令: %s", cmd.action)
	}
}

// writeLinesToFile 将行写入文件
func writeLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (c *SedCommand) Help() string {
	return `sed - 流编辑器（基础版）

用法:
  sed [选项] 'script' [文件...]
  command | sed 'script'

说明:
  对文本进行流式编辑处理。支持基本的sed命令。

选项:
  -i, --in-place      原地编辑文件
  -n, --quiet         静默模式（不自动打印）
  -e, --expression    指定sed表达式

支持的命令:
  s/pattern/replacement/[flags]   替换
    - 标志 g: 全局替换
    - 标志 i: 忽略大小写
  d                               删除行
  p                               打印行
  a\text                          在行后追加文本
  i\text                          在行前插入文本

示例:
  sed 's/old/new/' file.txt           # 替换第一个old为new
  sed 's/old/new/g' file.txt          # 替换所有old为new
  sed 's/old/new/gi' file.txt         # 忽略大小写全局替换
  sed '/pattern/d' file.txt           # 删除匹配的行
  sed -n '/pattern/p' file.txt        # 只打印匹配的行
  sed 's/foo/bar/g' -i file.txt       # 原地编辑
  echo "hello" | sed 's/hello/hi/'    # 管道使用`
}

func (c *SedCommand) ShortHelp() string {
	return "流编辑器（替换、删除等）"
}

