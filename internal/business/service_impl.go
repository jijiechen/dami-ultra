package business

import (
	"context"

	"github.com/google/uuid"
	"github.com/yuchanns/kong-exercise-microservices/internal/repos"
	"github.com/yuchanns/kong-exercise-microservices/internal/repos/models"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
)

type ImplService struct {
	repo repos.ServiceRepo
}

func NewImplService(repo repos.ServiceRepo) *ImplService {
	return &ImplService{repo}
}

// Create implements IService.
func (i *ImplService) Create(ctx context.Context, item *ServiceVersion) (err error) {
	service := &models.KonnectService{
		Name:        item.Service.Name,
		Description: item.Service.Description,
	}
	return helpers.Transaction(ctx, func(txCtx context.Context) (err error) {
		err = i.repo.Create(ctx, service)
		if err != nil {
			return
		}
		return i.repo.CreateVersion(ctx, &models.KonnectServiceVersion{
			KonnectServiceID: service.ID,
			Version:          uuid.NewString(),
			Host:             item.Version.Host,
			Port:             item.Version.Port,
			Path:             item.Version.Path,
			Protocol:         item.Version.Protocol,
		})
	})
}

// Get implements IService.
func (i *ImplService) Get(ctx context.Context, id int64) (item ServiceDetail, err error) {
	service, err := i.repo.Get(ctx, id)
	if err != nil {
		return
	}
	versions, err := i.repo.ListVersions(ctx, []int64{service.ID})
	if err != nil {
		return
	}
	item = ServiceDetail{
		Service: Service{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
		},
		Versions: make([]*Version, 0, len(versions)),
	}
	for _, version := range versions {
		item.Versions = append(item.Versions, &Version{
			Host:     version.Host,
			Port:     version.Port,
			Path:     version.Path,
			Protocol: version.Protocol,
		})
	}
	return
}

// List implements IService.
func (i *ImplService) List(ctx context.Context, opt ListOpt) (res ListRes, err error) {
	services, total, err := i.repo.
		List(ctx, opt.Keyword, (opt.Pager.Page-1)*opt.Pager.PageSize, opt.Pager.PageSize)
	if err != nil {
		return
	}
	res.Total = total
	if len(services) == 0 {
		return
	}
	serviceIDs := make([]int64, 0, len(services))
	for _, service := range services {
		serviceIDs = append(serviceIDs, service.ID)
	}
	versions, err := i.repo.ListVersions(ctx, serviceIDs)
	if err != nil {
		return
	}
	counter := make(map[int64]int32, len(versions))
	for _, version := range versions {
		counter[version.KonnectServiceID]++
	}
	for _, service := range services {
		res.Items = append(res.Items, &ServiceItem{
			Service: Service{
				ID:          service.ID,
				Name:        service.Name,
				Description: service.Description,
			},
			VersionCnt: counter[service.ID],
		})
	}
	return
}

// Update implements IService.
func (i *ImplService) Update(ctx context.Context, item *ServiceVersion) (err error) {
	service, err := i.repo.Get(ctx, item.Service.ID)
	if err != nil {
		return
	}
	return helpers.Transaction(ctx, func(ctx context.Context) (err error) {
		err = i.repo.Update(ctx, &models.KonnectService{
			ID:          item.Service.ID,
			Name:        item.Service.Name,
			Description: item.Service.Description,
		})
		if err != nil {
			return
		}
		return i.repo.CreateVersion(ctx, &models.KonnectServiceVersion{
			KonnectServiceID: service.ID,
			Version:          uuid.NewString(),
			Host:             item.Version.Host,
			Port:             item.Version.Port,
			Path:             item.Version.Path,
			Protocol:         item.Version.Protocol,
		})
	})
}

var _ IService = (*ImplService)(nil)
