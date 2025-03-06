package repos

import (
	"context"

	"github.com/yuchanns/kong-exercise-microservices/internal/repos/models"
	"github.com/yuchanns/kong-exercise-microservices/utils/helpers"
)

type ServiceRepo interface {
	List(ctx context.Context, keyword string, offset, limit int32) ([]*models.KonnectService, int64, error)
	ListVersions(ctx context.Context, serviceIDs []int64) ([]*models.KonnectServiceVersion, error)
	Create(ctx context.Context, service *models.KonnectService) error
	CreateVersion(ctx context.Context, service *models.KonnectServiceVersion) error
	Get(ctx context.Context, id int64) (models.KonnectService, error)
	Update(ctx context.Context, service *models.KonnectService) error
}

type ImplServiceRepo struct{}

// Update implements ServiceRepo.
func (i *ImplServiceRepo) Update(ctx context.Context, service *models.KonnectService) (err error) {
	return helpers.GetTenantDB(ctx).Model(&models.KonnectService{}).
		Where("id = ?", service.ID).
		Updates(map[string]interface{}{
			"name":        service.Name,
			"description": service.Description,
		}).Error
}

// Get implements ServiceRepo.
func (i *ImplServiceRepo) Get(ctx context.Context, id int64) (item models.KonnectService, err error) {
	err = helpers.GetTenantDB(ctx).Model(&models.KonnectService{}).
		Where("id = ?", id).
		First(&item).Error
	return
}

// Create implements ServiceRepo.
func (i *ImplServiceRepo) Create(ctx context.Context, service *models.KonnectService) error {
	return helpers.GetTenantDB(ctx).Create(service).Error
}

// CreateVersion implements ServiceRepo.
func (i *ImplServiceRepo) CreateVersion(ctx context.Context, service *models.KonnectServiceVersion) error {
	return helpers.GetTenantDB(ctx).Create(service).Error
}

// List implements ServiceRepo.
func (i *ImplServiceRepo) List(ctx context.Context, keyword string, offset int32, limit int32) (items []*models.KonnectService, total int64, err error) {
	q := helpers.GetTenantDB(ctx).Model(&models.KonnectService{})
	if keyword != "" {
		q.Where("name LIKE ?", "%"+keyword+"%")
	}
	err = q.Count(&total).Error
	if err != nil {
		return
	}
	err = q.Offset(int(offset)).Limit(int(limit)).
		Order("created_at DESC").Find(&items).Error
	return
}

// ListVersions implements ServiceRepo.
func (i *ImplServiceRepo) ListVersions(ctx context.Context, serviceIDs []int64) (items []*models.KonnectServiceVersion, err error) {
	err = helpers.GetTenantDB(ctx).Model(&models.KonnectServiceVersion{}).
		Where("konnect_service_id IN ?", serviceIDs).
		Order("created_at DESC").
		Find(&items).Error
	return
}

func NewImplServiceRepo() *ImplServiceRepo {
	return &ImplServiceRepo{}
}

var _ ServiceRepo = (*ImplServiceRepo)(nil)
