package services

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
)

type ProblemSetService interface {
	GetProblemSetListByEventId(ctx context.Context, eventId uuid.UUID) ([]entity.ProblemSet, error)
}

type problemSetService struct {
	problemSetRepo       repositories.ProblemSetRepository
	problemSetAssignRepo repositories.ProblemSetAssignRepository
}

func NewProblemSetService(problemSetRepo repositories.ProblemSetRepository, problemSetAssignRepo repositories.ProblemSetAssignRepository) ProblemSetService {
	return &problemSetService{problemSetRepo: problemSetRepo, problemSetAssignRepo: problemSetAssignRepo}
}

func (s *problemSetService) GetProblemSetListByEventId(ctx context.Context, eventId uuid.UUID) ([]entity.ProblemSet, error) {
	psas, err := s.problemSetAssignRepo.GetProblemSetAssignByEventID(ctx, eventId)
	if err != nil {
		return nil, err
	}
	var problemSets []entity.ProblemSet
	for _, psa := range psas {
		ps, err := s.problemSetRepo.GetProblemSetByID(ctx, psa.ProblemSetId)
		if err != nil {
			return nil, err
		}

		problemSets = append(problemSets, ps)
	}
	return problemSets, nil
}
