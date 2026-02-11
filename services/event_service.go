package services

import (
	"context"
	"errors"
	"strings"
	"time"

	dto "abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService interface {
	AuthorizeUserToEvent(ctx context.Context, slug string, accountId uuid.UUID) error
	List(ctx context.Context, accountId uuid.UUID, p entity.Pagination) ([]entity.Events, int64, error)
	DetailBySlug(ctx context.Context, slug string, accountId uuid.UUID) (dto.EventDetailResponse, error)
	JoinByCode(ctx context.Context, accountID uuid.UUID, code string) (dto.EventDetailResponse, error)
	GetStatus(ctx context.Context, slug string, accountId uuid.UUID) (eventStatus dto.EventStatus, err error)
	CreateEvent(ctx context.Context, req dto.CreateEventRequest) (entity.Events, error)
	UpdateEvent(ctx context.Context, id uuid.UUID, req dto.UpdateEventRequest) (entity.Events, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error
}

type eventService struct {
	paymentService  PaymentService
	eventsRepo      repositories.EventsRepository
	eventAssignRepo repositories.EventAssignRepository
}

func NewEventService(paymentService PaymentService, eventsRepo repositories.EventsRepository, eventAssignRepo repositories.EventAssignRepository) EventService {
	return &eventService{paymentService: paymentService, eventsRepo: eventsRepo, eventAssignRepo: eventAssignRepo}
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
func (s *eventService) List(ctx context.Context, accountId uuid.UUID, pagination entity.Pagination) ([]entity.Events, int64, error) {
	list, total, err := s.eventsRepo.ListVisible(ctx, accountId, &pagination)
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

	eventDetail := dto.EventDetailResponse{Data: &ev}

	if err != nil {
		return eventDetail, http_error.DATA_NOT_FOUND
	}

	assigned, err := s.eventAssignRepo.GetByEventAndAccount(ctx, ev.Id, accountID)

	if err == nil && assigned.Id != uuid.Nil {
		eventDetail.RegisterStatus = 1
		return eventDetail, http_error.ALREADY_REGISTERED_TO_EVENT
	}

	if err != nil && !errors.Is(err, http_error.DATA_NOT_FOUND) && !errors.Is(err, gorm.ErrRecordNotFound) {
		return eventDetail, err
	}

	if ev.Price != 0 {

		paymentEvent, err := s.paymentService.PayEvent(ctx, accountID, ev.Id, ev.Price)

		if err != nil {
			return eventDetail, err
		}

		eventDetail.EventPayment = paymentEvent

		if paymentEvent.Status != entity.PaymentStatusPaid {
			return eventDetail, http_error.PAYMENT_REQUIRED
		}

	}

	_, err = s.eventAssignRepo.Assign(ctx,
		entity.EventAssign{
			Id:        uuid.New(),
			AccountId: accountID,
			EventId:   ev.Id,
		},
	)

	if err != nil {
		return eventDetail, err
	}

	eventDetail.RegisterStatus = 1
	return eventDetail, err
}

func (s *eventService) GetStatus(ctx context.Context, slug string, accountId uuid.UUID) (eventStatus dto.EventStatus, err error) {

	event, err := s.DetailBySlug(ctx, slug, accountId)
	currentTime := time.Now()
	eventStatus.IsHasNotStarted = currentTime.Before(event.Data.StartEvent)
	eventStatus.IsFinished = currentTime.After(event.Data.EndEvent)
	eventStatus.IsOnGoing = !(eventStatus.IsHasNotStarted) && !(eventStatus.IsFinished)
	return eventStatus, err
}

func (s *eventService) CreateEvent(ctx context.Context, req dto.CreateEventRequest) (entity.Events, error) {
	startEvent, err := time.Parse(time.RFC3339, req.StartEvent)
	if err != nil {
		if startEvent.Before(time.Now()) {
			return entity.Events{}, http_error.EVENT_START_DATE_INVALID
		}
		return entity.Events{}, http_error.INVALID_DATE_FORMAT
	}
	endEvent, err := time.Parse(time.RFC3339, req.EndEvent)
	if err != nil {
		if endEvent.Before(startEvent) {
			return entity.Events{}, http_error.EVENT_END_DATE_INVALID
		}
		return entity.Events{}, http_error.INVALID_DATE_FORMAT
	}

	if err := utils.ValidateCode(req.EventCode); err != nil {
		return entity.Events{}, err
	}

	slugVal := strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
	if _, err := s.eventsRepo.GetBySlug(ctx, slugVal); err == nil {
		return entity.Events{}, http_error.DUPLICATE_DATA
	}
	if _, err := s.eventsRepo.GetByCode(ctx, req.EventCode); err == nil {
		return entity.Events{}, http_error.DUPLICATE_DATA
	}

	ev := entity.Events{
		Id:         uuid.New(),
		Title:      req.Title,
		Slug:       slugVal,
		StartEvent: startEvent,
		EndEvent:   endEvent,
		Overview:   req.Overview,
		ImgBanner:  req.ImgBanner,
		EventCode:  req.EventCode,
		IsPublic:   req.IsPublic,
	}
	return s.eventsRepo.Create(ctx, ev)
}

func (s *eventService) UpdateEvent(ctx context.Context, id uuid.UUID, req dto.UpdateEventRequest) (entity.Events, error) {
	existing, err := s.eventsRepo.GetByID(ctx, id)
	if err != nil {
		return entity.Events{}, http_error.DATA_NOT_FOUND
	}

	if req.Title != "" {
		existing.Title = req.Title
	}

	if req.StartEvent != "" {
		start, err := time.Parse(time.RFC3339, req.StartEvent)
		if err == nil {
			if start.Before(time.Now()) {
				return entity.Events{}, http_error.EVENT_START_DATE_INVALID
			}
			existing.StartEvent = start
		}
	}

	if req.EndEvent != "" {
		end, err := time.Parse(time.RFC3339, req.EndEvent)
		if err == nil {
			if end.Before(existing.StartEvent) {
				return entity.Events{}, http_error.EVENT_END_DATE_INVALID
			}
			existing.EndEvent = end
		}
	}

	if req.Overview != "" {
		existing.Overview = req.Overview
	}
	if req.ImgBanner != "" {
		existing.ImgBanner = req.ImgBanner
	}
	if req.IsPublic != nil {
		existing.IsPublic = *req.IsPublic
	}

	return s.eventsRepo.Update(ctx, existing)
}

func (s *eventService) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if _, err := s.eventsRepo.GetByID(ctx, id); err != nil {
		return http_error.DATA_NOT_FOUND
	}
	return s.eventsRepo.Delete(ctx, id)
}
