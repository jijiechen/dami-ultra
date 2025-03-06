package models

type KonnectService struct {
	BaseModel

	ID          int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description string `json:"description" gorm:"column:description;type:varchar(255);not null"`
	CreatedAt   string `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   string `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
}
