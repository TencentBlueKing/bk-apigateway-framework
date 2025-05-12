# bk-apigateway framework

## 说明

bk-apigateway framework 是一个基于 Django Rest Framework + drf-spectacular 的开发框架，用于快速开发 API 接口，部署在蓝鲸 PaaS 开发者中心，并接入到蓝鲸 API 网关。

能极大地简化开发者对接 API 网关的工作。

并且本身是一个 Python 项目，可以通过编码的方式实现接口的组合、协议转换、接口编排等功能。

## 特性

1. 封装了蓝鲸 PaaS 开发者中心的相关配置，开发者只需要关心 API 的实现以及声明，无需关心部署运行时
2. 集成了 drf-spectacular，开发者使用 `@extend_schema` 注解进行接口的声明，支持 OpenAPI 3.0 规范，自动生成接入蓝鲸 API 网关所需的 definition.yaml 和 resources.yaml
3. 封装了蓝鲸 API 网关的注册流程，开发者只需要在蓝鲸 PaaS 开发者中心进行发布，即可自动注册到蓝鲸 API 网关
4. 封装了蓝鲸 API 网关的调用流程，通过 `apigw_manager.drf.authentication.ApiGatewayJWTAuthentication`和 `apigw_manager.drf.permission.ApiGatewayPermission`实现解析网关请求中的 jwt，并且根据接口配置校验应用或用户。

## 应用场景

1. 网络协议转换，可以将网关的 HTTP 请求转换为其他协议的请求，例如 grpc/thrift
2. 接口组合与编排，可以将多个接口组合成一个接口，也可以在代码中做一些接口编排，例如调用接口 A，根据返回结果调用接口 B 或接口 C
3. 接口协议转换，例如将遗留系统的旧协议封装成新的协议，提供给调用方，或者可以将某些字段做映射/校验/类型转换/配置默认值等

## 开发步骤

1. 在蓝鲸 PaaS 开发者中心新建 API 网关插件，会基于网关插件模板初始化一个项目，项目中包含了部署到蓝鲸 PaaS 开发者中心的相关配置/接入到蓝鲸 API 网关的相关配置。并且包含一个 Demo 示例;
2. 开发者根据需求，开发相关的 API，并且使用 drf-spectacular 的 `@extend_schema` 注解进行接口的声明，可以在本地开发环境 swagger ui 查看接口配置渲染的效果，并且确认是否正确，也可以使用 django command 生成 definition.yaml 和 resources.yaml 进行验证;
3. 将代码提交到仓库，并且到蓝鲸 PaaS 开发者中心 - 插件开发进行发布;
4. 发布时，会自动触发构建，将服务部署到蓝鲸 PaaS 开发者中心，并且接入到蓝鲸 API 网关，开发者可以在蓝鲸 API 网关中查看接口的配置，以及调试接口;

## 本地开发

设置环境变量 (可以在项目跟路径新建一个`.envrc`文件，将下面内容放入文件中，启动时会自动加载；也可以在启动命令行终端中手动执行下面的内容)

```bash
export DEBUG=True
export IS_LOCAL=True
export BK_APIGW_NAME="demo"
export BK_API_URL_TMPL=http://bkapi.example.com/api/{api_name}/
export BKPAAS_APP_ID="demo"
export BKPAAS_APP_SECRET=358622d8-d3e7-4522-8f16-b5530776bbb8
export BKPAAS_DEFAULT_PREALLOCATED_URLS='{"dev": "http://0.0.0.0:8080/"}'
export BKPAAS_ENVIRONMENT=dev
export BKPAAS_PROCESS_TYPE=web
```

之后启动命令

```bash
python manage.py runserver 0.0.0.0:8080
```

可以访问 swagger ui 地址：http://0.0.0.0:8080/api/schema/swagger-ui/#/open
注意这个地址可以查看所有接口的文档，确认正确性，但是如果想要调试，需要将 settings.py 中的 REST_FRAMEWORK DEFAULT_AUTHENTICATION_CLASSES/DEFAULT_PERMISSION_CLASSES 注解掉

此时，日志文件在项目上层目录

```bash
# 将app_code换成应用名称
tail -f ../logs/{app_code}/*.log
```

配置完之后，可以本地生成 definition.yaml 和 resources.yaml 进行测试

```bash
python manage.py generate_definition_yaml && cat definition.yaml
python manage.py generate_resources_yaml && cat resources.yaml
```

## 注意

### 1. @extend_schema 只能放在 get/post/put/delete/patch 方法上，不能放在其他方法上

例如 drf 默认的 RetrieveAPIView，如果继承了这个基类，但是覆写的是`retrieve`，此时配置的 `@extend_schema`生成会有问题，例如 `parameters` 会丢失; 覆写 `get` 方法即可

具体可以在本地完成开发后，访问 swagger ui 确认接口相关的配置是否正确渲染

```python
class RetrieveAPIView(mixins.RetrieveModelMixin,
                      GenericAPIView):
    """
    Concrete view for retrieving a model instance.
    """
    def get(self, request, *args, **kwargs):
        return self.retrieve(request, *args, **kwargs)
```

此时接口实现继承了 `RetrieveAPIView`，必须实现`get()`而不是`retrieve()`, 并且`@extend_schema` 一定要配置在 `get` 方法上

```python
class DemoRetrieveApi(generics.RetrieveAPIView):
    ......
    @extend_schema(
        ......
    )
    def get(self, request, id, *args, **kwargs):
        ......
```

### 2. View 实现中需要配置 `app_verified_required` 和 `user_verified_required`

接口认证 ApiGatewayPermission 需要用到这两个属性，并且这两个属性用于当前 View 下所有接口的声明，注册到网关时，会根据这两个属性生成网关侧相应的配置。

具体参考 api/v1/views.py 中的示例

### 3. 发布后到蓝鲸 PaaS 开发者中心后，如何排查一些问题？

框架中关键路径上都有打印日志，日志级别是 debug, 所以可以通过修改日志级别来获取日志，排查问题。

可以设置环境变量 `BKAPP_LOG_LEVEL=DEBUG`, 之后重新发布。然后复现问题，在开发者中心日志中查看日志。
