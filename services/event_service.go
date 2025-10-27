package services

import (
	"context"
	"time"

	dto "abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
)

type EventService interface {
	List(ctx context.Context, p repositories.Pagination) ([]entity.Events, int64, error)
	DetailBySlug(ctx context.Context, slug string, accountId uuid.UUID) (dto.EventDetailResponse, error)
	JoinByCode(ctx context.Context, accountID uuid.UUID, code string) (dto.EventDetailResponse, error)
}

type eventService struct {
	eventsRepo      repositories.EventsRepository
	eventAssignRepo repositories.EventAssignRepository
}

func NewEventService(eventsRepo repositories.EventsRepository, eventAssignRepo repositories.EventAssignRepository) EventService {
	return &eventService{eventsRepo: eventsRepo, eventAssignRepo: eventAssignRepo}
}

func (s *eventService) List(ctx context.Context, pagination repositories.Pagination) ([]entity.Events, int64, error) {
	list, total, err := s.eventsRepo.GetAllPaginate(ctx, pagination)
	return list, total, err
}

func (s *eventService) DetailBySlug(ctx context.Context, slug string, accountId uuid.UUID) (dto.EventDetailResponse, error) {

	ev, err := s.eventsRepo.GetBySlug(ctx, slug)
	if err != nil {
		return dto.EventDetailResponse{}, err
	}

	assign, _ := s.eventAssignRepo.GetByEventAndAccount(ctx, ev.Id, accountId)
	if !ev.IsPublic && assign.Id == uuid.Nil {
		return dto.EventDetailResponse{}, http_error.UNAUTHORIZED
	}

	status := 0
	if assign.Id != uuid.Nil {
		status = 1
	}

	return dto.EventDetailResponse{Data: &ev, RegisterStatus: status}, nil
}

func (s *eventService) JoinByCode(ctx context.Context, accountID uuid.UUID, code string) (dto.EventDetailResponse, error) {
	ev, err := s.eventsRepo.GetByCode(ctx, code)
	if err != nil {
		return dto.EventDetailResponse{}, err
	}

	exist, err := s.eventAssignRepo.GetByEventAndAccount(ctx, ev.Id, accountID)
	if err == nil && exist.Id != uuid.Nil {
		return dto.EventDetailResponse{}, http_error.ALREADY_REGISTERED_TO_EVENT
	}

	assign := entity.EventAssign{EventId: ev.Id, AccountId: accountID, AssignedAt: time.Now()}
	if _, err := s.eventAssignRepo.Assign(ctx, assign); err != nil {
		return dto.EventDetailResponse{}, err
	}

	return dto.EventDetailResponse{Data: &ev, RegisterStatus: 1}, nil
}
