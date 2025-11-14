package services

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"
	repo "abdanhafidz.com/go-clean-layered-architecture/repositories"
)

type RegionService interface {
	SeedProvinces(ctx context.Context, provinces []entity.RegionProvince) error
	SeedCities(ctx context.Context, cities []entity.RegionCity) error
	ListProvinces(ctx context.Context) ([]entity.RegionProvince, error)
	ListCitiesByProvince(ctx context.Context, provinceId uint) ([]entity.RegionCity, error)
}

type regionService struct{ regionRepo repositories.RegionRepository }

func NewRegionService(regionRepo repo.RegionRepository) RegionService {
	return &regionService{regionRepo: regionRepo}
}

func (s *regionService) SeedProvinces(ctx context.Context, provinces []entity.RegionProvince) error {
	_, err := s.regionRepo.BulkCreateProvinces(ctx, provinces)
	return err
}

func (s *regionService) SeedCities(ctx context.Context, cities []entity.RegionCity) error {
	_, err := s.regionRepo.BulkCreateCities(ctx, cities)
	return err
}
func (s *regionService) ListProvinces(ctx context.Context) ([]entity.RegionProvince, error) {
	return s.regionRepo.ListProvinces(ctx)
}
func (s *regionService) ListCitiesByProvince(ctx context.Context, provinceId uint) ([]entity.RegionCity, error) {
	return s.regionRepo.ListCitiesByProvinceId(ctx, provinceId)
}
