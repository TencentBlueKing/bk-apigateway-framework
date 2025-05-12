// Package serializer ...
package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/infras/database"
	"bk.tencent.com/{{cookiecutter.project_name}}/{{cookiecutter.project_name}}/pkg/model"
)

// CategoryListRequest List Categories API 输入结构
type CategoryListRequest struct {
	Keyword string `form:"keyword" binding:"omitempty"`
}

// CategoryListResponse List Categories API 返回结构
type CategoryListResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Updater   string `json:"updater"`
	UpdatedAt string `json:"updatedAt"`
}

// CategoryCreateRequest Create Category API 输入结构
type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

// Validate ...
func (req *CategoryCreateRequest) Validate(c *gin.Context) error {
	tx := database.Client(c.Request.Context()).Where("name = ?", req.Name).First(&model.Category{})
	if tx.Error == nil {
		return errors.Errorf("category name `%s` already used", req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}

// CategoryCreateResponse Create Category API 输出结构
type CategoryCreateResponse struct {
	ID int64 `json:"id"`
}

// CategoryRetrieveResponse Retrieve Category API 返回结构
type CategoryRetrieveResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	Updater   string `json:"updater"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CategoryUpdateRequest Update Category API 输入结构
type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

// Validate ...
func (req *CategoryUpdateRequest) Validate(c *gin.Context) error {
	tx := database.Client(c.Request.Context()).
		Not("id = ?", c.Param("id")).
		Where("name = ?", req.Name).
		First(&model.Category{})
	if tx.Error == nil {
		return errors.Errorf("category name `%s` already used", req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}
