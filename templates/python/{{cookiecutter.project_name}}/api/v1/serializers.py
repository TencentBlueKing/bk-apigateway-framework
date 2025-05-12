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
