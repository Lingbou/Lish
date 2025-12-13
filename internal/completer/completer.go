package completer

import (
	"os"
	"path/filepath"
	"strings"
)

// Completer 自动补全器
type Completer struct {
	commands []string
}

// NewCompleter 创建补全器
func NewCompleter(commands []string) *Completer {
	return &Completer{
		commands: commands,
	}
}

// Do 实现 readline.AutoCompleter 接口
func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])

	// 分析输入
	parts := strings.Fields(lineStr)

	if len(parts) == 0 {
		// 空行，补全命令
		return c.completeCommand(""), 0
	}

	// 获取最后一个词
	lastWord := ""
	if len(lineStr) > 0 && lineStr[len(lineStr)-1] != ' ' {
		lastWord = parts[len(parts)-1]
	}

	// 判断是补全命令还是路径
	if len(parts) == 1 && lastWord == parts[0] {
		// 第一个词，补全命令
		return c.completeCommand(lastWord), len(lastWord)
	}

	// 补全文件路径
	return c.completePath(lastWord), len(lastWord)
}

// completeCommand 补全命令
func (c *Completer) completeCommand(prefix string) [][]rune {
	var matches [][]rune

	for _, cmd := range c.commands {
		if strings.HasPrefix(cmd, prefix) {
			matches = append(matches, []rune(cmd[len(prefix):]))
		}
	}

	return matches
}

// completePath 补全文件路径
func (c *Completer) completePath(prefix string) [][]rune {
	// 处理空路径
	if prefix == "" {
		prefix = "."
	}

	// 展开 ~ 为用户主目录
	if strings.HasPrefix(prefix, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			prefix = filepath.Join(home, prefix[1:])
		}
	}

	// 分离目录和文件名前缀
	dir := filepath.Dir(prefix)
	filePrefix := filepath.Base(prefix)

	// 如果以 / 结尾，说明已经是完整目录
	if strings.HasSuffix(prefix, string(filepath.Separator)) {
		dir = prefix
		filePrefix = ""
	}

	// 读取目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	// 查找匹配的文件
	var matches [][]rune
	basePrefixLen := len(filePrefix)

	for _, entry := range entries {
		name := entry.Name()

		// 跳过隐藏文件（除非用户明确输入了 .）
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(filePrefix, ".") {
			continue
		}

		// 检查前缀匹配
		if strings.HasPrefix(name, filePrefix) {
			suffix := name[basePrefixLen:]

			// 如果是目录，添加路径分隔符
			if entry.IsDir() {
				suffix += string(filepath.Separator)
			}

			matches = append(matches, []rune(suffix))
		}
	}

	return matches
}
