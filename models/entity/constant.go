package models

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

const (
	PaymentStatusPending = "PENDING"
	PaymentStatusPaid    = "PAID"
	PaymentStatusFailed  = "FAILED"
	PaymentStatusExpired = "EXPIRED"
)

const MB = 1024 * 1024

type Pagination struct {
	Limit          int
	Offset         int
	Search         string
	SortBy         string
	Order          string
	RegisterStatus *int
	Status         *string
}
