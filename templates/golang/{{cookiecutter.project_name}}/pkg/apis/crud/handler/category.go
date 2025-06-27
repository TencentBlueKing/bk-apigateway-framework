// Package handler ...
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/apis/crud/serializer"
	"bk.tencent.com/{{cookiecutter.project_name}}/pkg/infras/database"
	model2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/model"
	ginx2 "bk.tencent.com/{{cookiecutter.project_name}}/pkg/utils/ginx"
)

// ListCategories ...
//
//	@Summary	获取分类列表
//	@Tags		crud
//	@Success	200	{object}	ginx.Response{data=[]serializer.CategoryListResponse}
//	@Router		/api/categories [get]
func ListCategories(c *gin.Context) {
	var req serializer.CategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	tx := database.Client(c.Request.Context()).Model(&model2.Category{})
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		tx = tx.Where("LOWER(name) LIKE ?", keyword).Or("LOWER(updater) LIKE ?", keyword)
	}

	var categories []model2.Category
	if err := tx.Find(&categories).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.CategoryListResponse{}
	for _, category := range categories {
		respData = append(respData, serializer.CategoryListResponse{
			ID:        category.ID,
			Name:      category.Name,
			Updater:   category.Updater,
			UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
		})
	}
	ginx2.SetResp(c, http.StatusOK, respData)
}

// CreateCategory ...
//
//	@Summary	创建分类
//	@Tags		crud
//	@Param		body	body		serializer.CategoryCreateRequest	true	"创建分类请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/categories [post]
func CreateCategory(c *gin.Context) {
	var req serializer.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	category := model2.Category{
		Name: req.Name,
		BaseModel: model2.BaseModel{
			Creator: ginx2.GetUserID(c),
			Updater: ginx2.GetUserID(c),
		},
	}
	if err := database.Client(c.Request.Context()).Create(&category).Error; err != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginx2.SetResp(c, http.StatusCreated, serializer.CategoryCreateResponse{ID: category.ID})
}

// RetrieveCategory ...
//
//	@Summary	获取单个分类
//	@ID			get_category_by_id
//	@Tags		crud
//	@Param		id	path		int	true	"分类 ID"
//	@Success	200	{object}	ginx.Response{data=serializer.CategoryRetrieveResponse}
//	@Router		/api/categories/{id} [get]
func RetrieveCategory(c *gin.Context) {
	var category model2.Category

	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).First(&category)
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	respData := serializer.CategoryRetrieveResponse{
		ID:        category.ID,
		Name:      category.Name,
		Creator:   category.Creator,
		Updater:   category.Updater,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}
	ginx2.SetResp(c, http.StatusOK, respData)
}

// UpdateCategory ...
//
//	@Summary	更新分类
//	@ID			update_category_by_id
//	@Tags		crud
//	@Param		id		path	int									true	"分类 ID"
//	@Param		body	body	serializer.CategoryUpdateRequest	true	"更新分类请求体"
//	@Success	204		"No Content"
//	@Router		/api/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var req serializer.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx2.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	var category model2.Category
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", c.Param("id")).First(&category)
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	// 更新 DB 模型字段
	category.Name = req.Name
	category.Updater = ginx2.GetUserID(c)
	tx = database.Client(ctx).Save(&category)
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	ginx2.SetResp(c, http.StatusNoContent, nil)
}

// DestroyCategory ...
//
//	@Summary	删除分类
//	@ID			delete_category_by_id
//	@Tags		crud
//	@Param		id	path	int	true	"分类 ID"
//	@Success	204	"No Content"
//	@Router		/api/categories/{id} [delete]
func DestroyCategory(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&model2.Category{})
	if tx.Error != nil {
		ginx2.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx2.SetResp(c, http.StatusNoContent, nil)
}
