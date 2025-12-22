package models

import (
	"regexp"
)

const (
	StatusNotStarted = "NOT_STARTED"
	StatusInProgress = "IN_PROGRESS"
	StatusFinished   = "FINISHED"
)

const (
	EventStatusUpcoming = "UPCOMING"
	EventStatusOngoing  = "ONGOING"
	EventStatusEnded    = "ENDED"
)

const MB = 1024 * 1024

var CodeRegex = regexp.MustCompile(`^[a-zA-Z0-9]{6,12}$`)

type Pagination struct {
	Limit          int
	Offset         int
	Search         string
	SortBy         string
	Order          string
	RegisterStatus *int
	Status         *string
}
