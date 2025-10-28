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
	AuthorizeUserToEvent(ctx context.Context, slug string, accountId uuid.UUID) error
	List(ctx context.Context, p repositories.Pagination) ([]entity.Events, int64, error)
	DetailBySlug(ctx context.Context, slug string, accountId uuid.UUID) (dto.EventDetailResponse, error)
	JoinByCode(ctx context.Context, accountID uuid.UUID, code string) (dto.EventDetailResponse, error)
	QuizListByEvent(ctx context.Context, slug string) ([]entity.ProblemSet, error)
}

type eventService struct {
	problemSetService ProblemSetService
	eventsRepo        repositories.EventsRepository
	eventAssignRepo   repositories.EventAssignRepository
}

func NewEventService(eventsRepo repositories.EventsRepository, eventAssignRepo repositories.EventAssignRepository) EventService {
	return &eventService{eventsRepo: eventsRepo, eventAssignRepo: eventAssignRepo}
}

func (s *eventService) AuthorizeUserToEvent(ctx context.Context, slug string, accountId uuid.UUID) error {
	ev, err := s.eventsRepo.GetBySlug(ctx, slug)

	if err != nil {
		return err
	}

	evAssign, err := s.eventAssignRepo.GetByEventAndAccount(ctx, ev.Id, accountId)
	if err == nil && evAssign.Id != uuid.Nil {
		return nil
	}

	event, err := s.eventsRepo.GetByID(ctx, evAssign.EventId)

	if event.IsPublic {
		return http_error.NOT_REGISTERED_TO_EVENT
	} else if !event.IsPublic {
		return http_error.UNAUTHORIZED
	}

	return err
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

func (s *eventService) QuizListByEvent(ctx context.Context, slug string) ([]entity.ProblemSet, error) {
	ev, err := s.eventsRepo.GetBySlug(ctx, slug)
	psList, err := s.problemSetService.GetProblemSetListByEventId(ctx, ev.Id)
	return psList, err
}
