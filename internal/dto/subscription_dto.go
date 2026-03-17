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

// CreateSubscriptionRequest represents the request body for creating a subscription
// @Description Request body for creating a new subscription
type CreateSubscriptionRequest struct {
	Validator
	ServiceName string  `json:"service_name" validate:"required,min=1,max=100" example:"Netflix" enums:"Netflix,Spotify,Apple Music"`
	Price       int     `json:"price" validate:"required,gt=0" example:"999" minimum:"1"`
	UserID      string  `json:"user_id" validate:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   string  `json:"start_date" validate:"required,datetime=01-2006" example:"01-2024"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"  example:"12-2024"`
}

// UpdateSubscriptionRequest represents the request body for updating a subscription
// @Description Request body for updating an existing subscription
type UpdateSubscriptionRequest struct {
	Validator
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=100" example:"Netflix"`
	Price       *int    `json:"price,omitempty" validate:"omitempty,gt=0" example:"999"`
	EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006" example:"01-2024"`
}

// SumSubscriptionRequest represents the request body for collecting statistics
// @Description Request body for collecting subscription statistics
type SumSubscriptionRequest struct {
	Validator
	ServiceName *string `json:"service_name,omitempty" validate:"omitempty,min=1,max=100" example:"Netflix"`
	UserID      *string `json:"user_id,omitempty" validate:"omitempty,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartPeriod string  `json:"start_period" validate:"required,datetime=01-2006" example:"01-2024"`
	EndPeriod   string  `json:"end_period" validate:"required,datetime=01-2006" example:"12-2024"`
}

// SubscriptionResponse represents the response for a subscription
// @Description Response containing subscription details
type SumSubscriptionResponse struct {
	TotalCount int `json:"total_count" example:"10"`
	TotalSum   int `json:"total_sum" example:"15"`
}

// UpdateSubscriptionResponse represents the response for an updated subscription
// @Description Response containing updated subscription fields
type SubscriptionResponse struct {
	ID          uint       `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"Netflix"`
	Price       uint       `json:"price" example:"999"`
	UserID      uuid.UUID  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	StartDate   time.Time  `json:"start_date" example:"01-2024"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"12-2024"`
}

// SumSubscriptionResponse represents the statistics response
// @Description Response containing subscription statistics
type UpdateSubscriptionResponse struct {
	ID          uint       `json:"id" example:"1"`
	ServiceName *string    `json:"service_name,omitempty" example:"Netflix"`
	Price       *uint      `json:"price,omitempty" example:"999"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"12-2024"`
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

type ErrorResponse struct {
	Error string `json:"error"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}
