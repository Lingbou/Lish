package commands

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type CurlCommand struct {
	stdout *os.File
}

func NewCurlCommand(stdout *os.File) *CurlCommand {
	return &CurlCommand{stdout: stdout}
}

func (c *CurlCommand) Name() string {
	return "curl"
}

func (c *CurlCommand) Execute(ctx context.Context, args []string) error {
	flags := pflag.NewFlagSet("curl", pflag.ContinueOnError)
	method := flags.StringP("request", "X", "GET", "HTTP 方法")
	output := flags.StringP("output", "o", "", "保存到文件")
	headers := flags.StringSliceP("header", "H", []string{}, "自定义请求头")
	timeout := flags.IntP("timeout", "t", 30, "超时时间（秒）")

	if err := flags.Parse(args); err != nil {
		return err
	}

	urls := flags.Args()
	if len(urls) == 0 {
		return fmt.Errorf("curl: 需要指定 URL")
	}

	url := urls[0]

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(*timeout) * time.Second,
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, *method, url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 添加自定义请求头
	for _, header := range *headers {
		// 解析 "Key: Value" 格式
		var key, value string
		if _, err := fmt.Sscanf(header, "%s %s", &key, &value); err == nil {
			key = key[:len(key)-1] // 去掉冒号
			req.Header.Set(key, value)
		}
	}

	// 设置 User-Agent
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Lish-Curl/0.4.0")
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var writer io.Writer = c.stdout

	// 如果指定了输出文件
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			return fmt.Errorf("创建文件失败: %w", err)
		}
		defer file.Close()
		writer = file
	}

	// 显示状态码
	fmt.Fprintf(os.Stderr, "HTTP/%d.%d %s\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status)

	// 复制响应体
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	return nil
}

func (c *CurlCommand) Help() string {
	return `curl - HTTP 请求工具

用法:
  curl [选项] URL

选项:
  -X, --request     HTTP 方法（GET, POST, PUT, DELETE 等）
  -o, --output      保存响应到文件
  -H, --header      添加自定义请求头
  -t, --timeout     超时时间（秒，默认 30）

描述:
  发送 HTTP 请求并显示响应。

示例:
  curl https://api.github.com                    # GET 请求
  curl -X POST https://httpbin.org/post          # POST 请求
  curl -o page.html https://example.com          # 保存到文件
  curl -H "Accept: application/json" URL         # 自定义请求头
  curl -t 10 https://slow-site.com               # 10 秒超时`
}

func (c *CurlCommand) ShortHelp() string {
	return "HTTP 请求工具"
}
