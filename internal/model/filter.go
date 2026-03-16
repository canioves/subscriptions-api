package model

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionFilter struct {
	ServiceName *string
	UserID      *uuid.UUID
	StartPeriod time.Time
	EndPeriod   time.Time
}

type SubscriptionStat struct {
	TotalCount int `db:"total_count"`
	TotalSum   int `db:"total_sum"`
}
