package model

// Category 分类
type Category struct {
	BaseModel
	ID      int64   `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name" gorm:"type:varchar(32);unique;not null"`
	Entries []Entry `json:"entries" gorm:"foreignKey:CategoryID"`
}
