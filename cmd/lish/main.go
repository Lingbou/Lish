package main

import (
	"fmt"
	"os"

	"github.com/Lingbou/Lish/internal/shell"
)

func main() {
	// Create Shell
	sh, err := shell.NewShell()
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化 Shell 失败: %v\n", err)
		os.Exit(1)
	}

	// Init Shell
	if err := sh.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "初始化 Shell 失败: %v\n", err)
		os.Exit(1)
	}

	// Run Shell
	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "运行 Shell 失败: %v\n", err)
		os.Exit(1)
	}
}
