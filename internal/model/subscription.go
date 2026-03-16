package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uint       `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       uint       `db:"price"`
	UserID      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

type UpdateSubscription struct {
	ID          uint       `db:"id"`
	ServiceName *string    `db:"service_name"`
	Price       *uint      `db:"price"`
	EndDate     *time.Time `db:"end_date"`
}
