package services

import (
	"context"

	dto "abdanhafidz.com/go-clean-layered-architecture/models/dto"
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"
	"github.com/gosimple/slug"
)

type OptionService interface {
	CreateBulk(ctx context.Context, req []dto.OptionsRequest) error
	GetBySlug(ctx context.Context, slug string) (entity.Options, error)
}

type optionService struct{ optionRepo repositories.OptionRepository }

func NewOptionService(optionRepo repositories.OptionRepository) OptionService {
	return &optionService{optionRepo: optionRepo}
}

func (s *optionService) CreateBulk(ctx context.Context, req []dto.OptionsRequest) error {
	for _, item := range req {
		optionSlug := slug.Make(item.OptionName)
		cat, err := s.optionRepo.CreateOptionCategory(ctx, entity.OptionCategory{
			OptionName: item.OptionName,
			OptionSlug: optionSlug,
		})
		if err != nil {
			return err
		}

		for _, v := range item.OptionValue {
			_, err := s.optionRepo.CreateOptionValue(ctx, entity.OptionValues{
				OptionCategoryId: cat.Id,
				OptionValue:      v,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *optionService) GetBySlug(ctx context.Context, slug string) (entity.Options, error) {
	cat, err := s.optionRepo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		return entity.Options{}, err
	}
	vals, err := s.optionRepo.ListValuesByCategoryId(ctx, cat.Id)
	if err != nil {
		return entity.Options{}, err
	}
	// convert
	ov := make([]entity.OptionValues, 0, len(vals))
	for _, v := range vals {
		ov = append(ov, entity.OptionValues(v))
	}
	return entity.Options{OptionCategory: entity.OptionCategory(cat), OptionValues: ov}, nil
}
