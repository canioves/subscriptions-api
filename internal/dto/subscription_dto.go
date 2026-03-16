package dto

import (
	"fmt"
	"subscriptions-api/internal/model"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type Validator struct{}

type CreateSubscriptionRequest struct {
	Validator
	ServiceName string  `json:"service_name" validate:"required,min=1,max=100"`
	Price       int     `json:"price" validate:"required,gt=0"`
	UserID      string  `json:"user_id" validate:"required,uuid"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type UpdateSubscriptionRequest struct {
	Validator
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=100"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,gt=0"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
}

type SumSubscriptionRequest struct {
	Validator
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=100"`
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid"`
	StartPeriod string  `json:"start_period" validate:"required,datetime=01-2006"`
	EndPeriod   string  `json:"end_period" validate:"required,datetime=01-2006"`
}

type SumSubscriptionResponse struct {
	TotalCount int `json:"total_count"`
	TotalSum   int `json:"total_sum"`
}

type SubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       uint       `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *uint      `json:"price,omitempty"`
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

func ToUpdateSubscriptionResponse(updateSub *model.UpdateSubscription) *UpdateSubscriptionResponse {
	if updateSub == nil {
		return nil
	}

	return &UpdateSubscriptionResponse{
		ID:          updateSub.ID,
		ServiceName: updateSub.ServiceName,
		Price:       updateSub.Price,
		EndDate:     updateSub.EndDate,
	}
}

func ToSumSubscriptionResponse(stats *model.SubscriptionStat) *SumSubscriptionResponse {
	if stats == nil {
		return nil
	}

	return &SumSubscriptionResponse{
		TotalCount: stats.TotalCount,
		TotalSum:   stats.TotalSum,
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

func (req *UpdateSubscriptionRequest) IsEmpty() bool {
	return req.ServiceName == nil && req.Price == nil && req.EndDate == nil
}

func (Validator) validate(req any) error {
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

func (req *CreateSubscriptionRequest) Validate() error {
	return req.validate(req)
}

func (req *UpdateSubscriptionRequest) Validate() error {
	return req.validate(req)
}

func (req *SumSubscriptionRequest) Validate() error {
	return req.validate(req)
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
