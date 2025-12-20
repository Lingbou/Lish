package progress

import (
	"testing"
	"time"
)

func TestNewBar(t *testing.T) {
	bar := NewBar(100, "Test")
	
	if bar.Total != 100 {
		t.Errorf("Total = %v, want 100", bar.Total)
	}
	if bar.Current != 0 {
		t.Errorf("Current = %v, want 0", bar.Current)
	}
	if bar.Description != "Test" {
		t.Errorf("Description = %v, want Test", bar.Description)
	}
}

func TestBarUpdate(t *testing.T) {
	bar := NewBar(100, "Test")
	
	bar.Update(50)
	if bar.Current != 50 {
		t.Errorf("Current = %v, want 50", bar.Current)
	}
	
	bar.Update(30)
	if bar.Current != 80 {
		t.Errorf("Current = %v, want 80", bar.Current)
	}
}

func TestBarSet(t *testing.T) {
	bar := NewBar(100, "Test")
	
	bar.Set(75)
	if bar.Current != 75 {
		t.Errorf("Current = %v, want 75", bar.Current)
	}
}

func TestBarFinish(t *testing.T) {
	bar := NewBar(100, "Test")
	
	bar.Update(50)
	bar.Finish()
	
	if bar.Current != 100 {
		t.Errorf("Current = %v, want 100 after Finish", bar.Current)
	}
	if !bar.finished {
		t.Error("finished should be true after Finish")
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0B"},
		{100, "100B"},
		{1024, "1.0KB"},
		{1024 * 1024, "1.0MB"},
		{1024 * 1024 * 1024, "1.0GB"},
		{1536, "1.5KB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatSize(%d) = %v, want %v", tt.size, result, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		contains string
	}{
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m30s"},
		{3600 * time.Second, "1h0m0s"},
	}

	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result == "" {
				t.Error("formatDuration should not return empty string")
			}
		})
	}
}

func BenchmarkBarUpdate(b *testing.B) {
	bar := NewBar(1000000, "Benchmark")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bar.Update(1)
	}
}

func BenchmarkFormatSize(b *testing.B) {
	size := int64(1024 * 1024 * 512) // 512 MB
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatSize(size)
	}
}

