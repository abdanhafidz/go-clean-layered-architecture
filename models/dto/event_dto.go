package dto

import (
	entity "abdanhafidz.com/go-boilerplate/models/entity"
)

type EventDetailResponse struct {
	Data           *entity.Events
	RegisterStatus int `json:"register_status" binding:"required"`
}

type JoinEventRequest struct {
	EventCode string `json:"event_code" binding:"required"`
}

type EventStatus struct {
	IsHasNotStarted bool
	IsOnGoing       bool
	IsFinished      bool
}
