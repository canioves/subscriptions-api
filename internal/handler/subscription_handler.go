package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"subscriptions-api/internal/dto"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error while decoding request body", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := h.service.CreateSubscription(r.Context(), sub); err != nil {
		http.Error(w, "error while creating subscription", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := dto.ToSubscriptionResponse(sub)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error while decoding response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
