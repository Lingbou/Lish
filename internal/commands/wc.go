package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type WcCommand struct {
	stdout *os.File
}

func NewWcCommand(stdout *os.File) *WcCommand {
	return &WcCommand{stdout: stdout}
}

func (c *WcCommand) Name() string {
	return "wc"
}

func (c *WcCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("wc", pflag.ContinueOnError)
	countLines := flags.BoolP("lines", "l", false, "只统计行数")
	countWords := flags.BoolP("words", "w", false, "只统计单词数")
	countBytes := flags.BoolP("bytes", "c", false, "只统计字节数")
	
	if err := flags.Parse(args); err != nil {
		return err
	}
	
	files := flags.Args()
	if len(files) == 0 {
		return fmt.Errorf("wc: 需要指定文件")
	}
	
	// 如果没有指定选项，显示所有统计
	showAll := !*countLines && !*countWords && !*countBytes
	
	var totalLines, totalWords, totalBytes int64
	
	for _, filename := range files {
		lines, words, bytes, err := c.countFile(filename)
		if err != nil {
			return err
		}
		
		c.printStats(filename, lines, words, bytes, showAll, *countLines, *countWords, *countBytes)
		
		totalLines += lines
		totalWords += words
		totalBytes += bytes
	}
	
	// 如果有多个文件，显示总计
	if len(files) > 1 {
		c.printStats("total", totalLines, totalWords, totalBytes, showAll, *countLines, *countWords, *countBytes)
	}
	
	return nil
}

func (c *WcCommand) countFile(filename string) (lines, words, bytes int64, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines++
		bytes += int64(len(line)) + 1 // +1 for newline
		
		// 统计单词数
		fields := strings.Fields(line)
		words += int64(len(fields))
	}
	
	return lines, words, bytes, scanner.Err()
}

func (c *WcCommand) printStats(name string, lines, words, bytes int64, showAll, showLines, showWords, showBytes bool) {
	var output string
	
	if showAll || showLines {
		output += fmt.Sprintf("%8d", lines)
	}
	if showAll || showWords {
		output += fmt.Sprintf("%8d", words)
	}
	if showAll || showBytes {
		output += fmt.Sprintf("%8d", bytes)
	}
	
	output += fmt.Sprintf(" %s\n", name)
	fmt.Fprint(c.stdout, output)
}

func (c *WcCommand) Help() string {
	return `wc - 统计文件的行数、单词数和字节数

用法:
  wc [选项] 文件...

选项:
  -l, --lines  只显示行数
  -w, --words  只显示单词数
  -c, --bytes  只显示字节数

描述:
  统计文件的行数、单词数和字节数。
  如果不指定选项，显示所有统计信息。

示例:
  wc file.txt             # 显示所有统计
  wc -l file.txt          # 只显示行数
  wc -w file.txt          # 只显示单词数
  wc *.txt                # 统计多个文件并显示总计`
}

func (c *WcCommand) ShortHelp() string {
	return "统计文件行数/字数"
}

