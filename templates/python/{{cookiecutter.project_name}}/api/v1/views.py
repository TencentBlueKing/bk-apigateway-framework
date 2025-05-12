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
        operation_id="v1_demo",
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
                build_bk_header_rewrite(set={"X-Foo": "scope-resource"}, remove=["X-Bar"]),
                build_bk_ip_restriction(blacklist=["127.0.0.1", "127.0.0.2"]),
                build_bk_rate_limit(
                    default_period=60,
                    default_tokens=1000,
                    specific_app_limits=[("demo1", 3600, 1000)],
                ),
            ],
            # 匹配所有子路径，默认为 False
            match_subpath=False,
        ),
    )
    def get(self, request, id, *args, **kwargs):
        logger.debug(
            "request.jwt: %s, request.app: %s, request.user: %s",
            request.jwt and request.jwt.payload,
            request.app and request.app.bk_app_code,
            request.user and request.user.username,
        )
        slz = serializers.DemoRetrieveInputSLZ(data=request.query_params)
        slz.is_valid(raise_exception=True)

        data = slz.validated_data

        result = self.get_serializer(
            data={"message": f"Hello, {data['name']}! and my id is {id}"},
        )
        result.is_valid(raise_exception=True)
        return Response(result.data)
