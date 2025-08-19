package utils

import (
	"path/filepath"
	"testing"
)

func TestGetParentDir(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// 基础路径测试
		{
			name:     "绝对路径三级目录",
			input:    "/var/log/app",
			expected: "/var/log",
		},
		{
			name:     "绝对路径末尾带斜杠",
			input:    "/home/user/",
			expected: "/home",
		},
		{
			name:     "相对路径三级目录",
			input:    "src/pkg/utils",
			expected: "src/pkg",
		},
		// 特殊符号测试
		{
			name:     "包含点号路径",
			input:    "/opt/./app/current",
			expected: "/opt/app",
		},
		{
			name:     "包含双点号路径",
			input:    "/etc/nginx/../conf.d",
			expected: "/etc",
		},
		// 边界情况测试
		{
			name:     "根目录",
			input:    "/",
			expected: "/",
		},
		{
			name:     "单级目录",
			input:    "/tmp",
			expected: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行并验证结果
			got := GetParentDir(tt.input)
			if got != tt.expected {
				t.Errorf("\n输入: %s\n期望: %s\n实际: %s", tt.input, tt.expected, got)
			}
			// 附加验证：清理后的路径是否符合预期
			cleaned := filepath.Clean(tt.input)
			if parent := filepath.Dir(cleaned); parent != got {
				t.Errorf("内部验证失败！\n清理后路径: %s\nDir结果: %s\n函数结果: %s",
					cleaned, parent, got)
			}
		})
	}
}
