from django.urls import path

from . import views

urlpatterns = [
    path("demo/<int:id>/", views.DemoRetrieveApi.as_view(), name="api.v1.demo"),
]
