package dto

import (
	entity "abdanhafidz.com/go-boilerplate/models/entity"
)

type EventDetailResponse struct {
	Data           *entity.Events
	RegisterStatus int `json:"register_status"`
}

type JoinEventRequest struct {
	EventCode string `json:"event_code"`
}
