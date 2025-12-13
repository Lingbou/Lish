package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

// UniqCommand uniq 命令 - 去除或报告重复行
type UniqCommand struct{}

func (c *UniqCommand) Name() string {
	return "uniq"
}

func (c *UniqCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("uniq", flag.ContinueOnError)
	count := flags.BoolP("count", "c", false, "在每行前显示重复次数")
	repeated := flags.BoolP("repeated", "d", false, "只显示重复的行")
	uniqueOnly := flags.BoolP("unique", "u", false, "只显示唯一的行（不重复）")
	ignoreCase := flags.BoolP("ignore-case", "i", false, "忽略大小写")

	if err := flags.Parse(args); err != nil {
		return err
	}

	remaining := flags.Args()

	var lines []string
	var err error

	// 读取输入
	if len(remaining) == 0 {
		// 从标准输入读取
		lines, err = readLines(os.Stdin)
		if err != nil {
			return fmt.Errorf("读取输入失败: %w", err)
		}
	} else {
		// 从文件读取
		filename := remaining[0]
		lines, err = readLinesFromFile(filename)
		if err != nil {
			return fmt.Errorf("读取文件失败: %w", err)
		}
	}

	// 处理重复行
	result := processLines(lines, *count, *repeated, *uniqueOnly, *ignoreCase)

	// 输出结果
	for _, line := range result {
		fmt.Println(line)
	}

	return nil
}

// processLines 处理行，检测重复
func processLines(lines []string, showCount, showRepeated, showUnique, ignoreCase bool) []string {
	if len(lines) == 0 {
		return []string{}
	}

	type lineInfo struct {
		original string
		count    int
	}

	var result []string
	var current *lineInfo
	var processed []*lineInfo

	// 比较函数
	compare := func(a, b string) bool {
		if ignoreCase {
			return strings.ToLower(a) == strings.ToLower(b)
		}
		return a == b
	}

	// 统计相邻的重复行
	for _, line := range lines {
		if current == nil {
			current = &lineInfo{original: line, count: 1}
		} else if compare(current.original, line) {
			current.count++
		} else {
			processed = append(processed, current)
			current = &lineInfo{original: line, count: 1}
		}
	}
	if current != nil {
		processed = append(processed, current)
	}

	// 根据选项输出
	for _, info := range processed {
		// 过滤条件
		if showRepeated && info.count == 1 {
			continue // 只显示重复的，跳过唯一的
		}
		if showUnique && info.count > 1 {
			continue // 只显示唯一的，跳过重复的
		}

		// 输出格式
		if showCount {
			result = append(result, fmt.Sprintf("%7d %s", info.count, info.original))
		} else {
			result = append(result, info.original)
		}
	}

	return result
}

func (c *UniqCommand) Help() string {
	return `uniq - 去除或报告重复的相邻行

用法:
  uniq [选项] [文件]
  command | uniq [选项]

说明:
  过滤相邻的重复行。通常与 sort 一起使用。
  注意：uniq 只检测相邻的重复行，要去除所有重复行，
  请先使用 sort。

选项:
  -c, --count         在每行前显示重复次数
  -d, --repeated      只显示重复的行
  -u, --unique        只显示唯一的行（不重复的）
  -i, --ignore-case   比较时忽略大小写

示例:
  uniq file.txt                 # 去除相邻重复行
  uniq -c file.txt              # 显示每行出现次数
  uniq -d file.txt              # 只显示重复行
  uniq -u file.txt              # 只显示唯一行
  sort file.txt | uniq          # 排序后去重
  sort file.txt | uniq -c       # 统计每行出现次数
  cat file.txt | sort | uniq -c | sort -rn   # 按出现次数倒序`
}

func (c *UniqCommand) ShortHelp() string {
	return "去除或报告重复行"
}
