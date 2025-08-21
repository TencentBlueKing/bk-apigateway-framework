// Package handler ...
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TencentBlueKing/blueapps-go/pkg/utils/ginx"

	"github.com/TencentBlueKing/{{cookiecutter.project_name}}/pkg/apis/user/serializer"
)

// ListUsers ...
//  @ID		list_users
//	@Summary	查询user列表
//	@Tags		user
//	@Success	200	{object}	ginx.Response{data=[]serializer.UserListResponse}
//	@Router		/api/Users [get]
func ListUsers(c *gin.Context) {
	var req serializer.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	respData := []serializer.UserListResponse{}
	// TODO 查询数据
	ginx.SetResp(c, http.StatusOK, respData)
}

// CreateUser ...
//  @ID		create_user
//	@Summary	新增用户
//	@Tags		user
//	@Param		body	body		serializer.UserCreateRequest	true	"创建用户请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/Users [post]
func CreateUser(c *gin.Context) {
	var req serializer.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	// TODO 创建数据
	ginx.SetResp(c, http.StatusCreated, serializer.UserCreateResponse{ID: 1})
}
