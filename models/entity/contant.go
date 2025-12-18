package models

const (
	StatusNotStarted = "NOT_STARTED"
	StatusInProgress = "IN_PROGRESS"
	StatusFinished  = "FINISHED"
)

const MB = 1024 * 1024

type Pagination struct {
    Limit  int
    Offset int
    Search string
    SortBy string
    Order  string
}