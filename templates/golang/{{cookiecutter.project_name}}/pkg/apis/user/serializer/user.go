// Package serializer ...
package serializer

// UserListRequest List Users API 输入结构
type UserListRequest struct {
	Keyword string `form:"keyword" binding:"omitempty"`
}

// UserListResponse List Users API 返回结构
type UserListResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Updater   string `json:"updater"`
	UpdatedAt string `json:"updatedAt"`
}

// UserCreateRequest Create Users API 输入结构
type UserCreateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

// UserCreateResponse Create Users API 输出结构
type UserCreateResponse struct {
	ID int64 `json:"id"`
}
