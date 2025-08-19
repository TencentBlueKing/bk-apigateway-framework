package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/cmd"
)

// swagger api doc
//
//	@title			BkApp API DOC
//	@version		1.0.0
//	@description	蓝鲸开发者中心 SaaS 后台 API 文档
func main() {
	cmd.Execute()
}
