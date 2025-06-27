// Package async 提供一个简单的异步 / 定时任务封装：
// 1. 使用 cron 支持定时任务（cmd: scheduler）
// 2. 简单封装 goroutine 以支持异步任务
package async

import (
	"context"
	"reflect"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/async/task"
	log "bk.tencent.com/{{cookiecutter.project_name}}/pkg/logging"
)

// RegisteredTasks 已注册的任务
// 注意：任务函数最后一个返回值推荐为 error 类型
var RegisteredTasks = map[string]any{
	"CalcFib": task.CalcFib,
	// NOTE: SaaS 开发者可根据需求添加自定义任务
}

// ApplyTask 下发异步任务
func ApplyTask(ctx context.Context, name string, args []any) {
	go func() {
		taskFunc, ok := RegisteredTasks[name]
		if !ok {
			log.Errorf(ctx, "task func %s not found", name)
			return
		}

		taskArgs := []reflect.Value{reflect.ValueOf(ctx)}
		for _, arg := range args {
			taskArgs = append(taskArgs, reflect.ValueOf(arg))
		}
		values := reflect.ValueOf(taskFunc).Call(taskArgs)

		// 若任务执行有返回值，且最后一个返回值类型是 error 且不为 nil，需打印错误日志
		if length := len(values); length != 0 {
			if err, ok := values[length-1].Interface().(error); ok && err != nil {
				log.Errorf(ctx, "apply task %s with args %v error: %s", name, args, err)
			}
		}
	}()
}
