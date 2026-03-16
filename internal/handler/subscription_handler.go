package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subscriptions-api/internal/dto"
	appErrors "subscriptions-api/internal/errors"
	"subscriptions-api/internal/logger"
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
		http.Error(w, "Error while parsing ID parameter", http.StatusBadRequest)
		logger.Error("[HANDLER] Error while parsing ID parameter -> %s", err)
		return 0
	}
	if id < 0 {
		http.Error(w, "ID parameter must be greater than 0", http.StatusBadRequest)
		logger.Error("[HANDLER] ID parameter must be greater than 0")
		return 0
	}
	return uint(id)
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start creating subscription...")

	w.Header().Set("Content-Type", "application/json")

	var req dto.CreateSubscriptionRequest

	logger.Info("[HANDLER] Decoding request...")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		logger.Error("[HANDLER] Error while decoding request body -> %s", err)
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErr)
			for _, v := range validationErr.Errors {
				logger.Error("[HANDLER] Validation error -> %s", v)
			}
			return
		}
		http.Error(w, "Validation error", http.StatusBadRequest)
		logger.Error("[HANDLER] Validation error -> %s", err)
		return
	}

	parsedStartDate, err := dto.ParseDate(&req.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format: must be mm-yyyy", http.StatusBadRequest)
		logger.Error("[HANDLER] Invalid start_date format: must be mm-yyyy -> %s", err)
		return
	}
	parsedEndDate, err := dto.ParseDate(req.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format: must be mm-yyyy", http.StatusBadRequest)
		logger.Error("[HANDLER] Invalid end_date format: must be mm-yyyy -> %s", err)
		return
	}

	userUuid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		logger.Error("[HANDLER] Invalid user_id format -> %s", err)
		return
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       uint(req.Price),
		UserID:      userUuid,
		StartDate:   *parsedStartDate,
		EndDate:     parsedEndDate,
	}

	logger.Info("[HANDLER] Using service...")
	if err := h.service.CreateSubscription(r.Context(), sub); err != nil {
		http.Error(w, "Error while creating subscription", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while creating new subscription -> %s", err)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToSubscriptionResponse(sub)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error while decoding response", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while decoding response -> %s", err)
		return
	}
	logger.Info("[HANDLER] New subscription created successfully")
}

func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start receiving subscription...")
	w.Header().Set("Content-Type", "application/json")

	id := getIdParameter(w, r)

	logger.Info("[HANDLER] Using service...")
	sub, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		logger.Error("[HANDLER] Error while getting subscription -> %s", err)
		if errors.Is(err, appErrors.ErrNotFound) {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while getting subscription", http.StatusInternalServerError)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToSubscriptionResponse(sub)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error while decoding response", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while decoding response -> %s", err)
		return
	}
	logger.Info("[HANDLER] Subscription received successfully")
}

func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start receiving list of subscriptions...")
	w.Header().Set("Content-Type", "application/json")

	logger.Info("[HANDLER] Using service...")
	subs, err := h.service.ListSubscriptions(r.Context())
	if err != nil {
		http.Error(w, "Error while list subscriptions", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while list subscriptions -> %s", err)
		return
	}

	logger.Info("[HANDLER] Creating list of subscriptions response...")
	var responses []*dto.SubscriptionResponse
	for _, sub := range subs {
		response := dto.ToSubscriptionResponse(sub)
		responses = append(responses, response)
	}

	if err := json.NewEncoder(w).Encode(&responses); err != nil {
		http.Error(w, "Error while decoding response", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while decoding response -> %s", err)
		return
	}

	logger.Info("[HANDLER] List of subscriptions received successfully")
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start updating subscription...")
	w.Header().Set("Content-Type", "application/json")

	id := getIdParameter(w, r)

	logger.Info("[HANDLER] Decoding request...")
	var req dto.UpdateSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		logger.Error("[HANDLER] Error while decoding response -> %s", err)
		return
	}

	if req.IsEmpty() {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		logger.Info("[HANDLER] No fields to update")
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErr)
			for _, v := range validationErr.Errors {
				logger.Error("[HANDLER] Validation error -> %s", v)
			}
			return
		}
		http.Error(w, "Validation error", http.StatusBadRequest)
		logger.Error("[HANDLER] Validation error -> %s", err)
		return
	}

	var sub = &model.UpdateSubscription{}
	if req.ServiceName != nil {
		sub.ServiceName = req.ServiceName
	}
	if req.Price != nil {
		price := uint(*req.Price)
		sub.Price = &price

	}
	if req.EndDate != nil {
		parsedDate, err := dto.ParseDate(req.EndDate)
		if err != nil {
			http.Error(w, "Invalid end_date format: must be mm-yyyy", http.StatusBadRequest)
			logger.Error("[HANDLER] Invalid end_date format: must be mm-yyyy -> %s", err)
			return
		}
		sub.EndDate = parsedDate
	}

	logger.Info("[HANDLER] Using service...")
	if err := h.service.UpdateSubscription(r.Context(), id, sub); err != nil {
		logger.Error("[HANDLER] Error while updating subscription -> %s", err)
		if errors.Is(err, appErrors.ErrNotFound) {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error while updating subscription", http.StatusBadRequest)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToUpdateSubscriptionResponse(sub)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error while decoding response", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while decoding response -> %s", err)
		return
	}

	logger.Info("[HANDLER] Subscription updated successfully")
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start deleting subscription...")
	w.Header().Set("Content-Type", "application/json")

	id := getIdParameter(w, r)

	logger.Info("[HANDLER] Using service...")
	err := h.service.DeleteSubscription(r.Context(), id)
	if err != nil {
		http.Error(w, "Error while deleting subscription", http.StatusInternalServerError)
		if errors.Is(err, appErrors.ErrNotFound) {
			http.Error(w, "Subscription not found", http.StatusNotFound)
			return
		}
		logger.Error("[HANDLER] Error while deleting subscription -> %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	logger.Info("[HANDLER] Subscription deleted successfully")
}

func (h *SubscriptionHandler) CollectStats(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start collecting stats...")
	w.Header().Set("Content-Type", "application/json")

	logger.Info("[HANDLER] Decoding request...")
	var req dto.SumSubscriptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error while decoding request body", http.StatusBadRequest)
		logger.Error("[HANDLER] Error while decoding request body -> %s", err)
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErr)
			for _, v := range validationErr.Errors {
				logger.Error("[HANDLER] Validation error -> %s", v)
			}
			return
		}
		http.Error(w, "Validation error", http.StatusBadRequest)
		logger.Error("[HANDLER] Validation error -> %s", err)
		return
	}

	var filter = &model.SubscriptionFilter{}

	if req.ServiceName != nil {
		filter.ServiceName = req.ServiceName
	}
	if req.UserID != nil {
		userUuid, err := uuid.Parse(*req.UserID)
		if err != nil {
			http.Error(w, "Invalid user_id format", http.StatusBadRequest)
			logger.Error("[HANDLER] Invalid user_id format -> %s", err)
			return
		}
		filter.UserID = &userUuid
	}

	parsedStartPeriod, err := dto.ParseDate(&req.StartPeriod)
	if err != nil {
		http.Error(w, "Invalid end_date format: must be mm-yyyy", http.StatusBadRequest)
		logger.Error("[HANDLER] Invalid end_date format: must be mm-yyyy -> %s", err)
		return
	}
	filter.StartPeriod = *parsedStartPeriod

	parsedEndPeriod, err := dto.ParseDate(&req.EndPeriod)
	if err != nil {
		http.Error(w, "Invalid end_date format: must be mm-yyyy", http.StatusBadRequest)
		logger.Error("[HANDLER] Invalid end_date format: must be mm-yyyy -> %s", err)
		return
	}
	filter.EndPeriod = *parsedEndPeriod

	logger.Info("[HANDLER] Using service...")
	stats, err := h.service.CollectStats(r.Context(), filter)
	if err != nil {
		http.Error(w, "Error while collecting stats", http.StatusInternalServerError)
		logger.Error("[HANDLER] Error while collecting stats -> %s", err)
		return
	}

	logger.Info("[HANDLER] Creating stats response...")
	response := dto.ToSumSubscriptionResponse(stats)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "error while decoding response", http.StatusInternalServerError)
		logger.Error("[HANDLER] Invalid user_id format -> %s", err)
		return
	}

	logger.Info("[HANDLER] Subscriptions stats collected")
}
