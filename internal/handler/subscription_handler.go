package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"subscriptions-api/internal/dto"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func getIdParameter(w http.ResponseWriter, r *http.Request) uint {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "error while parsing id parameter", http.StatusBadRequest)
		log.Println(err)
		return 0
	}
	if id < 0 {
		http.Error(w, "id parameter must be greater than 0", http.StatusBadRequest)
		log.Println("id parameter must be greater than 0")
		return 0
	}
	return uint(id)
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "error while decoding request body", http.StatusBadRequest)
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

func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := getIdParameter(w, r)

	sub, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		http.Error(w, "error while getting subscription", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := dto.ToSubscriptionResponse(sub)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error while decoding response", http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	subs, err := h.service.ListSubscriptions(r.Context())
	if err != nil {
		http.Error(w, "error while getting all subscriptions", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var responses []*dto.SubscriptionResponse
	for _, sub := range subs {
		response := dto.ToSubscriptionResponse(sub)
		responses = append(responses, response)
	}

	if err := json.NewEncoder(w).Encode(&responses); err != nil {
		http.Error(w, "error while decoding response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := getIdParameter(w, r)

	err := h.service.UpdateSubscription(r.Context(), id, sub)
}
