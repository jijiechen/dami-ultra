package models

type KonnectServiceVersion struct {
	BaseModel

	ID               int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	KonnectServiceID int64  `json:"konnect_service_id" gorm:"column:konnect_service_id;type:int;not null"`
	Version          string `json:"version" gorm:"column:version;type:varchar(255);not null"`
	Host             string `json:"host" gorm:"column:host;type:varchar(255);not null"`
	Port             int    `json:"port" gorm:"column:port;type:int;not null"`
	Path             string `json:"path" gorm:"column:path;type:varchar(255);not null"`
	Protocol         string `json:"protocol" gorm:"column:protocol;type:varchar(255);not null"`
	CreatedAt        string `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        string `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
}
