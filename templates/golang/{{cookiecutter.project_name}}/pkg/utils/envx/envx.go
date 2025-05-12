// Package envx 提供环境变量相关工具
package envx

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Get 读取环境变量，支持默认值
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// MustGet 读取环境变量，若不存在则 panic
func MustGet(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("required environment variable %s unset", key))
}

// GetEnvList 读取环境变量列表，支持默认值
func GetEnvList(key string, defaultList []string) []string {
	// 检查环境变量是否存在
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultList
	}
	// 处理空字符串情况
	if val == "" {
		return []string{}
	}
	// 分割并清理空格
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		result = append(result, trimmed)
	}
	return result
}

// MustGetEnvJSON 解析 JSON 环境变量
func MustGetEnvJSON(key string) (map[string]string, error) {
	jsonStr := os.Getenv(key)
	if jsonStr == "" {
		return nil, fmt.Errorf("env: %s 未设置", key)
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}
	return result, nil
}
