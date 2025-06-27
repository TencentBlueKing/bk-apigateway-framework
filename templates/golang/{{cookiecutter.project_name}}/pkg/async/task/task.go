// Package task 包含异步任务实现
package task

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/model"
)

// Fibonacci 斐波那契数的递归实现，因为性能很差所以适合模拟需要长时间运行的后台任务
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// CalcFib 计算斐波那契数任务
func CalcFib(ctx context.Context, n float64) (int, error) {
	// 由于 json Unmarshal 会把整数 & 浮点数都解析为 float64 类型，这由任务处理类型转换
	nInt := int(n)

	task := model.Task{
		Name:      "CalcFib",
		Args:      []byte(fmt.Sprintf("{\"n\": %d}", nInt)),
		StartedAt: time.Now(),
	}
	if err := database.Client(ctx).Create(&task).Error; err != nil {
		return 0, err
	}

	// 执行计算任务
	fibN := fibonacci(nInt)

	// 回填执行结果
	task.Result = []byte(strconv.Itoa(fibN))
	task.Duration = time.Since(task.StartedAt)
	if err := database.Client(ctx).Save(&task).Error; err != nil {
		return 0, err
	}

	return fibN, nil
}
