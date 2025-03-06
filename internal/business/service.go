package business

import "context"

type ListOpt struct {
	*Pager

	Keyword string `json:"keyword" form:"keyword"`
}

type GetOpt struct {
	ID int64 `json:"id" form:"id" uri:"id" validate:"required"`
}

type ListRes struct {
	Total int64          `json:"total"`
	Items []*ServiceItem `json:"items"`
}

type ServiceItem struct {
	Service
	VersionCnt int32 `json:"version_cnt"`
}

type ServiceDetail struct {
	Service
	Versions []*Version `json:"versions"`
}

type Service struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type Version struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,gte=0,lte=65535"`
	Path     string `json:"path" validate:"required"`
	Protocol string `json:"protocol" validate:"required,oneof=http https grpc"`
}

type ServiceVersion struct {
	Service
	Version
}

type IService interface {
	List(ctx context.Context, opt ListOpt) (ListRes, error)
	Get(ctx context.Context, id int64) (ServiceDetail, error)
	Create(ctx context.Context, item *ServiceVersion) error
	Update(ctx context.Context, item *ServiceVersion) error
}
