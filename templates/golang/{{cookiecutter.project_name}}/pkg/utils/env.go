package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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

// GetEnvJSON 解析 JSON 环境变量
func GetEnvJSON(key string, defaultVal map[string]string) (map[string]string, error) {
	jsonStr := os.Getenv(key)
	if jsonStr == "" {
		return defaultVal, nil
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}
	return result, nil
}
