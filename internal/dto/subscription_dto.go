package dto

import (
	"fmt"
	"subscriptions-api/internal/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" validate:"required,min=1,max=100"`
	Price       int     `json:"price" validate:"required,gt=0"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
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

func ParseDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	pattern := "01-2006"
	parsedDate, err := time.Parse(pattern, *s)
	if err != nil {
		return nil, fmt.Errorf("failed to conver string to date: %w", err)
	}
	return &parsedDate, nil
}

func (req *CreateSubscriptionRequest) Validate() error {
	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)

		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			errors[field] = getErrorMessage(field, tag)
		}

		return &ValidationError{Errors: errors}
	}
	return nil
}

func getErrorMessage(field string, tag string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "min":
		return field + " is too short"
	case "max":
		return field + " is too long"
	case "gt":
		return field + " must be greater than 0"
	case "uuid":
		return field + " must be a valid UUID"
	case "datetime":
		return field + " must be in MM-YYYY format"
	default:
		return field + " is invalid"
	}
}

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
