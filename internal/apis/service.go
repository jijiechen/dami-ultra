package apis

import (
	"context"

	"github.com/yuchanns/kong-exercise-microservices/internal/business"
)

type Service struct {
	biz business.IService
}

func NewService(biz business.IService) *Service {
	return &Service{biz}
}

func (s *Service) List(ctx context.Context, r business.ListOpt) (business.ListRes, error) {
	r.Pager = business.GetPager(r.Pager)
	return s.biz.List(ctx, r)
}

func (s *Service) Get(ctx context.Context, r business.GetOpt) (business.ServiceDetail, error) {
	return s.biz.Get(ctx, r.ID)
}

func (s *Service) Create(ctx context.Context, r business.ServiceVersion) (any, error) {
	return nil, s.biz.Create(ctx, &r)
}

func (s *Service) Update(ctx context.Context, r business.ServiceVersion) (any, error) {
	return nil, s.biz.Update(ctx, &r)
}
