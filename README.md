# bk-apigateway-framework

![img](https://github.com/TencentBlueKing/blueking-apigateway/blob/master/docs/resource/img/blueking_apigateway_zh.png)
---

[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/bk-apigateway-framework/blob/main/LICENSE.txt) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/bk-apigateway-framework/pulls)

[English](README_EN.md)

## 简介

该仓库包含蓝鲸 API 网关-可编程网关开发框架， 提供了 Go/Python 两种语言的框架，开发者可以基于该框架构建自己的可编程网关。

## 总览

- Python 框架
  - [文档](./docs/python.md)
  - [框架模板](./templates/python/)
- Golang 框架
  - [文档](./docs/golang.md)
  - [框架模板](./templates/golang/)

## 使用方式

依赖于 [cookiecutter](https://github.com/cookiecutter/cookiecutter) 工具

```bash
pip install cookiecutter

cookiecutter https://github.com/TencentBlueKing/bk-apigateway-framework/ --directory templates/python

cookiecutter https://github.com/TencentBlueKing/bk-apigateway-framework/ --directory templates/golang
```

## 蓝鲸社区

- [BK-APIGateway](https://github.com/TencentBlueKing/blueking-apigateway): 蓝鲸API网关提供了高性能、高可用的 API 托管服务。
- [BK-CI](https://github.com/Tencent/bk-ci)：蓝鲸持续集成平台是一个开源的持续集成和持续交付系统，可以轻松将你的研发流程呈现到你面前。
- [BK-BCS](https://github.com/Tencent/bk-bcs)：蓝鲸容器管理平台是以容器技术为基础，为微服务业务提供编排管理的基础服务平台。
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：蓝鲸PaaS平台是一个开放式的开发平台，让开发者可以方便快捷地创建、开发、部署和管理SaaS应用。
- [BK-SOPS](https://github.com/Tencent/bk-sops)：标准运维（SOPS）是通过可视化的图形界面进行任务流程编排和执行的系统，是蓝鲸体系中一款轻量级的调度编排类SaaS产品。
- [BK-CMDB](https://github.com/Tencent/bk-cmdb)：蓝鲸配置平台是一个面向资产及应用的企业级配置管理平台。

## 贡献

如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests，为蓝鲸开源社区贡献力量。

## 协议

基于 MIT 协议，详细请参考[LICENSE](LICENSE.txt)