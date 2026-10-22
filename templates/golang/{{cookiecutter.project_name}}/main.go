package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/cmd"
	_ "github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/version" // 注入构建版本信息
)

// swagger api doc
//
//	@title			BkApp API DOC
//	@version		1.0.0
//	@description	蓝鲸开发者中心 SaaS 后台 API 文档
func main() {
	cmd.Execute()
}
