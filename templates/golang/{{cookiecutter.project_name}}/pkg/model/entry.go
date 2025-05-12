package model

// Entry 条目
type Entry struct {
	BaseModel
	CategoryID int64    `json:"categoryID" gorm:"not null"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	ID    int64   `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name" gorm:"type:varchar(64);unique;not null"`
	Desc  string  `json:"desc" gorm:"type:text;null"`
	Price float32 `json:"price" gorm:"not null"`
}
