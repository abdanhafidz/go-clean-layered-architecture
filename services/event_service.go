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
	List(ctx context.Context, accountId uuid.UUID, p repositories.Pagination) ([]entity.Events, int64, error)
	DetailBySlug(ctx context.Context, slug string, accountId uuid.UUID) (dto.EventDetailResponse, error)
	JoinByCode(ctx context.Context, accountID uuid.UUID, code string) (dto.EventDetailResponse, error)
	GetStatus(ctx context.Context, slug string, accountId uuid.UUID) (eventStatus dto.EventStatus, err error)
}

type eventService struct {
	eventsRepo      repositories.EventsRepository
	eventAssignRepo repositories.EventAssignRepository
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
func (s *eventService) List(ctx context.Context, accountId uuid.UUID, pagination repositories.Pagination) ([]entity.Events, int64, error) {
	evList := []entity.Events{}
	evPublicList, total, err := s.eventsRepo.ListPublic(ctx, &pagination)
	evList = append(evList, evPublicList...)

	if err != nil {
		return []entity.Events{}, 0, err
	}

	evAssignList, err := s.eventAssignRepo.ListByAccount(ctx, accountId)

	if err != nil {
		return []entity.Events{}, 0, err
	}

	for _, evAssign := range evAssignList {
		evPrivate, err := s.eventsRepo.GetByID(ctx, evAssign.EventId)
		if err != nil {
			return []entity.Events{}, 0, err
		}
		evList = append(evList, evPrivate)
	}

	return evList, total, err
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

func (s *eventService) GetStatus(ctx context.Context, slug string, accountId uuid.UUID) (eventStatus dto.EventStatus, err error) {
	
	
	event, err := s.DetailBySlug(ctx, slug, accountId)
	currentTime := time.Now()
	eventStatus.IsHasNotStarted = currentTime.Before(event.Data.StartEvent)
	eventStatus.IsFinished = currentTime.Before(event.Data.EndEvent)
	eventStatus.IsOnGoing = !(eventStatus.IsHasNotStarted) && !(eventStatus.IsFinished)

	return eventStatus, err


}
