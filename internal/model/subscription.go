package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint
	ServiceName string
	Price       uint
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}
