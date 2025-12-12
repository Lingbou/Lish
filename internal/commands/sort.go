package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// SortCommand sort 命令 - 排序文本行
type SortCommand struct{}

func (c *SortCommand) Name() string {
	return "sort"
}

func (c *SortCommand) Execute(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("sort", flag.ContinueOnError)
	reverse := flags.BoolP("reverse", "r", false, "反向排序")
	numeric := flags.BoolP("numeric-sort", "n", false, "数字排序")
	unique := flags.BoolP("unique", "u", false, "去除重复行")
	fieldNum := flags.IntP("key", "k", 0, "按指定字段排序（从1开始）")
	separator := flags.StringP("field-separator", "t", "", "字段分隔符")
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
		for _, filename := range remaining {
			fileLines, err := readLinesFromFile(filename)
			if err != nil {
				return fmt.Errorf("读取文件 %s 失败: %w", filename, err)
			}
			lines = append(lines, fileLines...)
		}
	}

	// 排序
	if *numeric {
		// 数字排序
		sort.Slice(lines, func(i, j int) bool {
			ni := extractNumber(lines[i], *fieldNum, *separator)
			nj := extractNumber(lines[j], *fieldNum, *separator)
			if *reverse {
				return ni > nj
			}
			return ni < nj
		})
	} else if *fieldNum > 0 {
		// 按字段排序
		sort.Slice(lines, func(i, j int) bool {
			fi := extractField(lines[i], *fieldNum, *separator)
			fj := extractField(lines[j], *fieldNum, *separator)
			
			if *ignoreCase {
				fi = strings.ToLower(fi)
				fj = strings.ToLower(fj)
			}
			
			if *reverse {
				return fi > fj
			}
			return fi < fj
		})
	} else {
		// 普通字符串排序
		sort.Slice(lines, func(i, j int) bool {
			li := lines[i]
			lj := lines[j]
			
			if *ignoreCase {
				li = strings.ToLower(li)
				lj = strings.ToLower(lj)
			}
			
			if *reverse {
				return li > lj
			}
			return li < lj
		})
	}

	// 去重（如果指定）
	if *unique {
		lines = removeDuplicates(lines)
	}

	// 输出结果
	for _, line := range lines {
		fmt.Println(line)
	}

	return nil
}

// extractField 提取指定字段
func extractField(line string, fieldNum int, separator string) string {
	if fieldNum <= 0 {
		return line
	}

	var fields []string
	if separator == "" {
		// 默认使用空白字符分隔
		fields = strings.Fields(line)
	} else {
		fields = strings.Split(line, separator)
	}

	if fieldNum > len(fields) {
		return ""
	}

	return fields[fieldNum-1]
}

// extractNumber 提取数字
func extractNumber(line string, fieldNum int, separator string) float64 {
	field := extractField(line, fieldNum, separator)
	if field == "" {
		field = line
	}

	// 尝试解析数字
	num, err := strconv.ParseFloat(strings.TrimSpace(field), 64)
	if err != nil {
		return 0
	}
	return num
}

// removeDuplicates 去除重复行
func removeDuplicates(lines []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return result
}

// readLines 从 Reader 读取所有行
func readLines(file *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// readLinesFromFile 从文件读取所有行
func readLinesFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readLines(file)
}

func (c *SortCommand) Help() string {
	return `sort - 排序文本行

用法:
  sort [选项] [文件...]
  command | sort [选项]

说明:
  对文本行进行排序。如果未指定文件，从标准输入读取。

选项:
  -r, --reverse           反向排序
  -n, --numeric-sort      按数字大小排序
  -u, --unique            去除重复行
  -k, --key <N>           按第 N 个字段排序（从1开始）
  -t, --field-separator   指定字段分隔符
  -i, --ignore-case       忽略大小写

示例:
  sort file.txt                    # 排序文件
  sort -r file.txt                 # 反向排序
  sort -n numbers.txt              # 数字排序
  sort -u file.txt                 # 去重排序
  sort -k 2 file.txt               # 按第2个字段排序
  sort -t: -k 3 /etc/passwd        # 使用:分隔，按第3字段排序
  ls | sort                        # 对ls输出排序
  cat file1.txt file2.txt | sort   # 合并排序`
}

func (c *SortCommand) ShortHelp() string {
	return "排序文本行"
}

