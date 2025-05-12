# bk-apigateway-framework

![img](https://github.com/TencentBlueKing/blueking-apigateway/blob/master/docs/resource/img/blueking_apigateway_en.png)
---

[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/bk-apigateway-framework/blob/main/LICENSE.txt) [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/bk-apigateway-framework/pulls)

[中文](README.md)

## Introduction

The bk-apigateway-framework is a framework for building programmable gateways using Go and Python. It provides a set of tools and libraries to help developers build their own programmable gateways.

## Overview

- Python Framework
  - [Documentation](./docs/python.md)
  - [Framework Template](./templates/python/)
- Golang Framework
  - [Documentation](./docs/golang.md)
  - [Framework Template](./templates/golang/)

## Usage

Dependencies: [cookiecutter](https://github.com/cookiecutter/cookiecutter)

```bash
pip install cookiecutter

cookiecutter https://github.com/TencentBlueKing/bk-apigateway-framework/ --directory templates/python

cookiecutter https://github.com/TencentBlueKing/bk-apigateway-framework/ --directory templates/golang
```

## BlueKing Community

- [BK-APIGateway](https://github.com/TencentBlueKing/blueking-apigateway): a high-performance and highly available API hosting service.
- [BK-CI](https://github.com/Tencent/bk-ci)：a continuous integration and continuous delivery system that can easily present your R & D process to you.
- [BK-BCS](https://github.com/Tencent/bk-bcs)：a basic container service platform which provides orchestration and management for micro-service business.
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：an development platform that allows developers to create, develop, deploy and manage SaaS applications easily and quickly.
- [BK-SOPS](https://github.com/Tencent/bk-sops)：an lightweight scheduling SaaS  for task flow scheduling and execution through a visual graphical interface.
- [BK-CMDB](https://github.com/Tencent/bk-cmdb)：an enterprise-level configuration management platform for assets and applications.

## Contributing

If you have good ideas or suggestions, please let us know by Issues or Pull Requests and contribute to the Blue Whale Open Source Community.

## License

Based on the MIT protocol, Please refer to [LICENSE](LICENSE.txt)