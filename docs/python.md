## 背景

bk-apigateway framework 是一个基于 [Django Rest Framework](https://www.django-rest-framework.org/) + [drf-spectacular](https://drf-spectacular.readthedocs.io/en/latest/) 的开发框架，用于快速开发 API 接口，部署在蓝鲸 PaaS 开发者中心，并接入到蓝鲸 API 网关。

能极大地简化开发者对接 API 网关的工作。

并且本身是一个 Python 项目，可以通过编码的方式实现接口的组合、协议转换、接口编排等功能。

## 特性

1. 封装了蓝鲸 PaaS 开发者中心的相关配置，开发者只需要关心 API 的实现以及声明，无需关心部署运行时
2. 集成了 [drf-spectacular](https://drf-spectacular.readthedocs.io/en/latest/)，开发者使用 `@extend_schema` 注解进行接口的声明，支持 OpenAPI 3.0 规范，自动生成接入蓝鲸 API 网关所需的 definition.yaml 和 resources.yaml
3. 封装了蓝鲸 API 网关的注册流程，开发者只需要在蓝鲸 PaaS 开发者中心进行发布，即可自动注册到蓝鲸 API 网关
4. 封装了蓝鲸 API 网关的调用流程，通过 `apigw_manager.drf.authentication.ApiGatewayJWTAuthentication`和 `apigw_manager.drf.permission.ApiGatewayPermission`实现解析网关请求中的 jwt，并且根据接口配置校验应用或用户。

## 使用场景

1. 网络协议转换，可以将网关的 HTTP 请求转换为其他协议的请求，例如 grpc/thrift
2. 接口组合与编排，可以将多个接口组合成一个接口，也可以在代码中做一些接口编排，例如调用接口 A，根据返回结果调用接口 B 或接口 C
3. 接口协议转换，例如将遗留系统的旧协议封装成新的协议，提供给调用方，或者可以将某些字段做映射/校验/类型转换/配置默认值等

## 使用

### 1. 开发步骤

1. 在蓝鲸 PaaS 开发者中心新建 API 网关插件，会基于网关插件模板初始化一个项目，项目中包含了部署到蓝鲸 PaaS 开发者中心的相关配置/接入到蓝鲸 API 网关的相关配置。并且包含一个 Demo 示例;
2. 开发者根据需求，开发相关的 API，并且使用 drf-spectacular 的 `@extend_schema` 注解进行接口的声明，可以在本地开发环境 swagger ui 查看接口配置渲染的效果，并且确认是否正确，也可以使用 django command 生成 definition.yaml 和 resources.yaml 进行验证;
3. 将代码提交到仓库，并且到蓝鲸 PaaS 开发者中心 - 插件开发进行发布;
4. 发布时，会自动触发构建，将服务部署到蓝鲸 PaaS 开发者中心，并且接入到蓝鲸 API 网关，开发者可以在蓝鲸 API 网关中查看接口的配置，以及调试接口;

### 2. 蓝鲸 API 网关新建可编程网关

新建成功后，可以下载到对应的代码

### 3. 本地开发

默认框架中带有示例代码，可以参考目录`api/v1/`下的 serializer 和 view 写法法

设置环境变量 (可以在项目跟路径新建一个`.envrc`文件，将下面内容放入文件中，启动时会自动加载；也可以在启动命令行终端中手动执行下面的内容)

```bash
export DEBUG=True
export IS_LOCAL=True
export BK_APIGW_NAME="demo"
export BK_API_URL_TMPL=http://bkapi.example.com/api/{api_name}/
export BKPAAS_APP_ID="demo"
export BKPAAS_APP_SECRET=358622d8-d3e7-4522-8f16-b5530776bbb8
export BKPAAS_DEFAULT_PREALLOCATED_URLS='{"prod": "http://0.0.0.0:8080/"}'
export BKPAAS_ENVIRONMENT=prod
export BKPAAS_PROCESS_TYPE=web
```

之后启动命令

```bash
python manage.py runserver 0.0.0.0:8080
```

可以访问 swagger ui 地址：http://0.0.0.0:8080/api/schema/swagger-ui/#/open
注意这个地址可以查看所有接口的文档，确认正确性，但是如果想要调试，需要将 settings.py 中的 `REST_FRAMEWORK DEFAULT_AUTHENTICATION_CLASSES/DEFAULT_PERMISSION_CLASSES` 注解掉

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

### 3. 提交代码并发布

确认发布结果

### 4. 蓝鲸 API 网关确认并进行在线调试

使用 蓝鲸 API 网关 的 在线调试 功能，可以进行在线调试

## 接口配置示例

### view

添加蓝鲸 API 网关模型相关配置，注意可以配置一系列开关及插件，需要根据接口需求处理

```python
import logging

from apigw_manager.drf.utils import gen_apigateway_resource_config
from apigw_manager.plugin.config import (
    build_bk_cors,
    build_bk_header_rewrite,
    build_bk_ip_restriction,
    build_bk_rate_limit,
)
from drf_spectacular.utils import extend_schema
from rest_framework import generics
from rest_framework.response import Response

from . import serializers

logger = logging.getLogger("app")


# 更多关于 drf 视图类的说明 (about drf Concrete View Classes)
# https://www.django-rest-framework.org/api-guide/generic-views/#concrete-view-classes
class DemoRetrieveApi(generics.RetrieveAPIView):
    serializer_class = serializers.DemoRetrieveOutputSLZ

    # 是否开启应用认证，对 APIView 下所有的方法生效 (用于 ApiGatewayPermission 中的校验)
    app_verified_required = True
    # 是否开启用户认证，对 APIView 下所有的方法生效 (用于 ApiGatewayPermission 中的校验)
    user_verified_required = True

    # 更多关于 `@extend_schema` 的说明 (more details about extend_schema in drf_spectacular)
    # https://drf-spectacular.readthedocs.io/en/latest/readme.html#customization-by-using-extend-schema
    @extend_schema(
        # 全局唯一，避免冲突
        operation_id="api_v1_demo",
        description="这是一个 demo api",
        parameters=[
            serializers.DemoRetrieveInputSLZ,
        ],
        responses={200: serializers.DemoRetrieveOutputSLZ},
        # 标签，用于同步时过滤掉不需要注册的接口，以及注册网关时资源对应打的标签
        tags=["open"],
        extensions=gen_apigateway_resource_config(
            # 是否公开，不公开在文档中心/应用申请网关权限资源列表中不可见
            is_public=True,
            # 是否允许申请权限，不允许的话在应用申请网关权限资源列表中不可见
            allow_apply_permission=True,
            # 是否开启用户认证，这里必须引用类变量 (因为需要保证网关侧配置调用 jwt 到 当前项目 permission_class 中的校验一致)
            user_verified_required=user_verified_required,
            # 是否开启应用认证，这里必须引用类变量 (因为需要保证网关侧配置调用 jwt 到 当前项目 permission_class 中的校验一致)
            app_verified_required=app_verified_required,
            # 是否校验资源权限，是的话将会校验应用是否有调用这个资源的权限，前置条件：开启应用认证
            resource_permission_required=True,
            description_en="this is a demo api",
            # 插件配置，类型为 List[Dict], 用于声明作用在这个资源上的插件，可以参考官方文档
            # 没有特殊需求的话默认不需要开启任何插件，如果需要开启插件，可以参考下面的例子
            # NOTE: 注意不要直接复制下面的内容到你的接口定义中，除非你知道每个插件配置后产生的影响
            plugin_configs=[
                build_bk_cors(),
                build_bk_header_rewrite(set={"X-Foo": "test"}, remove=["X-Bar"]),
                build_bk_ip_restriction(blacklist=["127.0.0.1", "192.168.2.1"]),
                build_bk_rate_limit(
                    default_period=60,
                    default_tokens=1000,
                    specific_app_limits=[("demo", 3600, 1000)],
                ),
            ],
            # 匹配所有子路径，默认为 False
            match_subpath=False,
        ),
    )
    def get(self, request, id, *args, **kwargs):
        ........
```

### serializer

添加 input/output 的 example

```python
from drf_spectacular.utils import OpenApiExample, extend_schema_serializer
from rest_framework import serializers


# more details about extend_schema_serializer in drf_spectacular
# https://drf-spectacular.readthedocs.io/en/latest/customization.html#step-4-extend-schema-serializer
@extend_schema_serializer(
    examples=[
        OpenApiExample(
            "example",
            value={
                "name": "world",
            },
            request_only=True,  # signal that example only applies to responses
        ),
    ]
)
class DemoRetrieveInputSLZ(serializers.Serializer):
    name = serializers.CharField(max_length=100)


# more details about extend_schema_serializer in drf_spectacular
# https://drf-spectacular.readthedocs.io/en/latest/customization.html#step-4-extend-schema-serializer
@extend_schema_serializer(
    exclude_fields=("type",),  # schema ignore these fields
    examples=[
        OpenApiExample(
            "example",
            value={
                "message": "hello, world, and my id is 1",
            },
            response_only=True,  # signal that example only applies to responses
        ),
    ],
)
class DemoRetrieveOutputSLZ(serializers.Serializer):
    message = serializers.CharField()
    type = serializers.CharField(required=False)
```

## 网关配置说明

所有网关/环境/版本相关的配置在 `config/settings.py`中，可以根据关键字定位并修改配置

### 1. 网关基本属性

> 可以通过环境变量注入 (建议), 也可以直接修改配置代码

```python
# 网关是否公开，公开则其他开发者可见/可申请权限
BK_APIGW_IS_PUBLIC = str(env.bool("BK_APIGW_IS_PUBLIC", default=True)).lower()
# if BK_APIGW_IS_OFFICIAL is True, the BK_APIGW_NAME should be start with `bk-`
BK_APIGW_IS_OFFICIAL = 1 if env.bool("BK_APIGW_IS_OFFICIAL", default=False) else 10
# 网关管理员，请将负责人加入列表中
BK_APIGW_MAINTAINERS = env.list("BK_APIGW_MAINTAINERS", default=["admin"])
```

### 2. 环境相关配置

超时时间，环境变量以及环境插件

```python
# 网关接口最大超时时间
BK_APIGW_STAG_BACKEND_TIMEOUT = 60

# 声明网关不同环境的环境变量
stag_env_vars = {
    "foo": "bar"
}
prod_env_vars = {
    # "foo": "bar"
}

# 声明网关不同环境的插件配置
# https://github.com/TencentBlueKing/bkpaas-python-sdk/blob/master/sdks/apigw-manager/docs/plugin-use-guide.md
# 注意，这里声明的插件配置会作用在对应环境的所有资源上，所以谨慎声明，确保你知道每个插件配置后产生的影响
stag_plugin_configs = build_stage_plugin_config_for_definition_yaml(
    [
        build_bk_cors(),
        build_bk_header_rewrite(set={"X-Foo": "scope-stage-stag"}, remove=["X-Bar"]),
        build_bk_ip_restriction(blacklist=["192.168.2.1", "192.168.2.2"]),
        build_bk_rate_limit(
            default_period=60,
            default_tokens=1000,
            specific_app_limits=[("demo3", 3600, 1000)],
        ),
    ]
)
prod_plugin_configs = build_stage_plugin_config_for_definition_yaml(
    [
        # build_bk_cors(),
        # build_bk_header_rewrite(set={"X-Foo": "scope-stage-prod"}, remove=["X-Bar"]),
        # build_bk_ip_restriction(blacklist=["192.168.1.1", "192.168.1.2"]),
        # build_bk_rate_limit(
        #     default_period=60,
        #     default_tokens=1000,
        #     specific_app_limits=[("demo2", 3600, 1000)],
        # ),
    ]
)
```

### 3. 主动授权

> 可以通过环境变量注入 (建议), 也可以直接修改配置代码

```python
# 主动授权，网关主动给应用，添加访问网关所有资源
BK_APIGW_GRANT_PERMISSION_DIMENSION_GATEWAY_APP_CODES = env.list(
    "BK_APIGW_GRANT_PERMISSION_DIMENSION_GATEWAY_APP_CODES", default=[]
)
BK_APIGW_GRANT_PERMISSION_DIMENSION_RESOURCE_APP_CODES = {
    # app_code: [resource_name1, resource_name2]
    "demo": ["v1_demo"],
}
```

### 4. 版本号及版本日志

> 可以通过环境变量注入 (建议), 也可以直接修改配置代码

```python
# release settings
# YOU CAN CHANGE THE RELEASE INFO, use env vars or just change the default below
# 1.0.0+stag or 1.0.0+prod
BK_APIGW_RELEASE_VERSION = (
    # NOTE: 每次部署必须强制版本号变更，否则代码变更版本号不变，不会打出新版本
    # log: resource_version 1.0.3+stag already exists, skip creating
    env.str("BK_APIGW_RELEASE_VERSION", default="1.0.0") + "+" + BK_APIGW_STAGE_NAME
)


BK_APIGW_RELEASE_TITLE = env.str("BK_APIGW_RELEASE_TITLE", default=f"gateway release(stage={BK_APIGW_STAGE_NAME})")
BK_APIGW_RELEASE_COMMENT = env.str(
    "BK_APIGW_RELEASE_COMMENT",
    default=f"auto release by bk-apigw-plugin-runtime(stage={BK_APIGW_STAGE_NAME})",
)
```

## 注意事项

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


### 4. 如何同步 `MCP Server` 相关配置
api 声明时通过 `enable_mcp=True` 开启 `MCP` 功能，并且注意确认请求参数。
```python
class DemoRetrieveApi(generics.RetrieveAPIView):
    ......
    @extend_schema(
     # 是否开启MCP
     enable_mcp=True,
     # 是否有请求参数，对于已经声明：parameters，requestBody参数的不用设置
     none_schema=True,
    )
    def get(self, request, id, *args, **kwargs):
        ......
```
配置：

```python
# 是否开启同步 MCP Server
BK_APIGW_STAGE_ENABLE_MCP_SERVERS = False
stage_mcp_servers = {
    "stag": [
        {
            "name": "mcp_server1",
            # 中文名
            "title": "示例 mcpserver1"
            "description": "mcp_server1",
            # 主动授权 app_code
            "target_app_codes": [APP_CODE],
            "labels": ["demo1"],
            # 是否启用：1-启用，0-停止
            "status": 1,
            # 协议类型：sse/streamable_http; 推荐使用 streamable_http 更稳定
            "protocol_type": "streamable_http",
            # 是否公开
            "is_public": True,
            # 添加的资源列表(如果不指定则会将符合规范的都加入)
            "tools": []
        }
    ],
    "prod": [
        {
            "name": "mcp_server1",
            "title": "示例 mcpserver1"
            "description": "mcp_server1",
            "target_app_codes": [APP_CODE],
            "labels": ["demo1"],
            "status": 0,
            "protocol_type": "streamable_http",
            "is_public": False,
            "tools": []
        },
        {
            "name": "mcp_server2",
            "title": "示例 mcpserver2"
            "description": "mcp_server2",
            "target_app_codes": [APP_CODE],
            "labels": ["demo2"],
            "status": 1,
            "protocol_type": "streamable_http",
            "is_public": True,
            "tools": ["demo2"]
        }
    ]
}

BK_APIGW_STAGE_MCP_SERVERS = stage_mcp_servers.get(bkpaas_environment, [])
```
