package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"gorm.io/gorm"
)

type RegionRepository interface {
	BulkCreateProvinces(ctx context.Context, provinces []entity.RegionProvince) ([]entity.RegionProvince, error)
	BulkCreateCities(ctx context.Context, cities []entity.RegionCity) ([]entity.RegionCity, error)
	ListProvinces(ctx context.Context) ([]entity.RegionProvince, error)
	ListCitiesByProvinceId(ctx context.Context, provinceId uint) ([]entity.RegionCity, error)
}

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{db: db}
}

func (r *regionRepository) BulkCreateProvinces(ctx context.Context, provinces []entity.RegionProvince) ([]entity.RegionProvince, error) {
	if len(provinces) == 0 {
		return provinces, nil
	}
	if err := r.db.WithContext(ctx).Create(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

func (r *regionRepository) BulkCreateCities(ctx context.Context, cities []entity.RegionCity) ([]entity.RegionCity, error) {
	if len(cities) == 0 {
		return cities, nil
	}
	if err := r.db.WithContext(ctx).Create(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *regionRepository) ListProvinces(ctx context.Context) ([]entity.RegionProvince, error) {
	var list []entity.RegionProvince
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *regionRepository) ListCitiesByProvinceId(ctx context.Context, provinceId uint) ([]entity.RegionCity, error) {
	var list []entity.RegionCity
	if err := r.db.WithContext(ctx).Where("province_id = ?", provinceId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
