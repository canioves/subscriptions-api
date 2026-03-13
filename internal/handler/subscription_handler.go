package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"subscriptions-api/internal/dto"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/service"

	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErr)
			return
		}
		http.Error(w, "validation error", http.StatusBadRequest)
		return
	}

	parsedStartDate, err := dto.ParseDate(&req.StartDate)
	if err != nil {
		http.Error(w, "invalid start_date format: must be mm-yyyy", http.StatusBadRequest)
		log.Println(err)
		return
	}
	parsedEndDate, err := dto.ParseDate(req.EndDate)
	if err != nil {
		http.Error(w, "invalid end_date format: must be mm-yyyy", http.StatusBadRequest)
		log.Println(err)
		return
	}

	userUuid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		log.Println(err)
		return
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       uint(req.Price),
		UserID:      userUuid,
		StartDate:   *parsedStartDate,
		EndDate:     parsedEndDate,
	}

	if err := h.service.CreateSubscription(r.Context(), sub); err != nil {
		http.Error(w, "error while creating subscription", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := dto.ToSubscriptionResponse(sub)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error while decoding response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
