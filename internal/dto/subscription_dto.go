package dto

import (
	"subscriptions-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	ServiceName string    `json:"service_name"`
	Price       uint      `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       uint       `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func ToSubscriptionResponse(sub *model.Subscription) *SubscriptionResponse {
	if sub == nil {
		return nil
	}

	return &SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
}

func (req *SubscriptionRequest) ParseDate() (*time.Time, error) {
	if req == nil {
		return nil, nil
	}

}
