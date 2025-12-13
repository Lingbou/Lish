package progress

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// Bar 进度条
type Bar struct {
	Total       int64
	Current     int64
	Width       int
	StartTime   time.Time
	Description string
	mu          sync.Mutex
	lastUpdate  time.Time
	finished    bool
}

// NewBar 创建新的进度条
func NewBar(total int64, description string) *Bar {
	return &Bar{
		Total:       total,
		Current:     0,
		Width:       50,
		StartTime:   time.Now(),
		Description: description,
		lastUpdate:  time.Now(),
		finished:    false,
	}
}

// Update 更新进度
func (b *Bar) Update(n int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Current += n

	// 限制更新频率（每100ms更新一次）
	now := time.Now()
	if now.Sub(b.lastUpdate) < 100*time.Millisecond && b.Current < b.Total {
		return
	}
	b.lastUpdate = now

	b.display()
}

// Set 设置当前进度
func (b *Bar) Set(current int64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Current = current
	b.display()
}

// Finish 完成进度条
func (b *Bar) Finish() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.finished {
		return
	}

	b.Current = b.Total
	b.finished = true
	b.display()
	fmt.Println() // 换行
}

// display 显示进度条
func (b *Bar) display() {
	if b.Total <= 0 {
		return
	}

	percent := float64(b.Current) / float64(b.Total) * 100
	if percent > 100 {
		percent = 100
	}

	// 计算进度条长度
	filled := int(float64(b.Width) * percent / 100)
	if filled > b.Width {
		filled = b.Width
	}

	// 构建进度条
	bar := strings.Repeat("█", filled) + strings.Repeat("░", b.Width-filled)

	// 计算速度和剩余时间
	elapsed := time.Since(b.StartTime)
	speed := float64(b.Current) / elapsed.Seconds()
	eta := ""
	if speed > 0 && b.Current < b.Total {
		remaining := time.Duration(float64(b.Total-b.Current)/speed) * time.Second
		eta = fmt.Sprintf(" ETA: %s", formatDuration(remaining))
	}

	// 格式化大小
	currentSize := formatSize(b.Current)
	totalSize := formatSize(b.Total)
	speedStr := formatSize(int64(speed)) + "/s"

	// 显示进度
	fmt.Printf("\r%s [%s] %.1f%% %s/%s %s%s",
		b.Description, bar, percent, currentSize, totalSize, speedStr, eta)
}

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}

	units := []string{"KB", "MB", "GB", "TB"}
	div := int64(unit)
	exp := 0

	for n := size / unit; n >= unit && exp < len(units)-1; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%s", float64(size)/float64(div), units[exp])
}

// formatDuration 格式化时间
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

// Reader 带进度的 io.Reader
type Reader struct {
	io.Reader
	bar *Bar
}

// NewReader 创建带进度的 Reader
func NewReader(r io.Reader, total int64, description string) *Reader {
	return &Reader{
		Reader: r,
		bar:    NewBar(total, description),
	}
}

// Read 读取数据并更新进度
func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.bar.Update(int64(n))
	if err == io.EOF {
		r.bar.Finish()
	}
	return
}

// Writer 带进度的 io.Writer
type Writer struct {
	io.Writer
	bar *Bar
}

// NewWriter 创建带进度的 Writer
func NewWriter(w io.Writer, total int64, description string) *Writer {
	return &Writer{
		Writer: w,
		bar:    NewBar(total, description),
	}
}

// Write 写入数据并更新进度
func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	w.bar.Update(int64(n))
	return
}

// Finish 完成写入
func (w *Writer) Finish() {
	w.bar.Finish()
}
