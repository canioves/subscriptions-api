package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"subscriptions-api/internal/dto"
	appErrors "subscriptions-api/internal/errors"
	"subscriptions-api/internal/logger"
	"subscriptions-api/internal/model"
	"subscriptions-api/internal/service"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

func getIdParameter(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, fmt.Errorf("ID parameter must be a number")
	}
	if id < 0 {
		return 0, fmt.Errorf("ID parameter must be greater than 0")
	}
	return uint(id), nil
}

func (h *SubscriptionHandler) sendError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
	if err != nil {
		logger.Error("[HANDLER] %s -> %s", message, err)
	}
}

func (h *SubscriptionHandler) sendValidationError(w http.ResponseWriter, validationErr *dto.ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(validationErr)
	for _, v := range validationErr.Errors {
		logger.Error("[HANDLER] Validation error -> %s", v)
	}
}

func (h *SubscriptionHandler) sendSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.sendError(w, http.StatusInternalServerError, "Error while encoding response", err)
		}
	}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "Subscription creation request"
// @Success 201 {object} dto.SubscriptionResponse "Subscription created successfully"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start creating subscription...")
	defer r.Body.Close()

	var req dto.CreateSubscriptionRequest

	logger.Info("[HANDLER] Decoding request...")
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Error while decoding request body", err)
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			h.sendValidationError(w, validationErr)
			return
		}
		h.sendError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	parsedStartDate, err := dto.ParseDate(&req.StartDate)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid start_date format: must be mm-yyyy", err)
		return
	}

	var parsedEndDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err = dto.ParseDate(req.EndDate)
		if err != nil {
			h.sendError(w, http.StatusBadRequest, "Invalid end_date format: must be mm-yyyy", err)
			return
		}
	}

	userUuid, err := uuid.Parse(req.UserID)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid user_id format", err)
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
		h.sendError(w, http.StatusInternalServerError, "Error while creating new subscription", err)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToSubscriptionResponse(sub)
	h.sendSuccess(w, http.StatusCreated, response)
	logger.Info("[HANDLER] New subscription created successfully")
}

// GetSubscription godoc
// @Summary Get a subscription by ID
// @Description Retrieve a subscription by its unique ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Success 200 {object} dto.SubscriptionResponse "Subscription found"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start receiving subscription...")

	id, err := getIdParameter(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	logger.Info("[HANDLER] Using service...")
	sub, err := h.service.GetSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			h.sendError(w, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Error while getting subscription", err)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToSubscriptionResponse(sub)
	h.sendSuccess(w, http.StatusOK, response)
	logger.Info("[HANDLER] Subscription received successfully")
}

// ListSubscriptions godoc
// @Summary List all subscriptions
// @Description Retrieve a list of all subscriptions
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} dto.SubscriptionResponse "List of subscriptions"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start receiving list of subscriptions...")

	logger.Info("[HANDLER] Using service...")
	subs, err := h.service.ListSubscriptions(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Error while listing subscriptions", err)
		return
	}

	logger.Info("[HANDLER] Creating list of subscriptions response...")
	var responses []*dto.SubscriptionResponse
	for _, sub := range subs {
		response := dto.ToSubscriptionResponse(sub)
		responses = append(responses, response)
	}

	h.sendSuccess(w, http.StatusOK, responses)
	logger.Info("[HANDLER] List of subscriptions received successfully")
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Param request body dto.UpdateSubscriptionRequest true "Subscription update request"
// @Success 200 {object} dto.UpdateSubscriptionResponse "Subscription updated successfully"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start updating subscription...")
	defer r.Body.Close()

	id, err := getIdParameter(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	logger.Info("[HANDLER] Decoding request...")
	var req dto.UpdateSubscriptionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Error while decoding request body", err)
		return
	}

	if req.IsEmpty() {
		h.sendError(w, http.StatusBadRequest, "No fields to update", nil)
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			h.sendValidationError(w, validationErr)
			return
		}
		h.sendError(w, http.StatusBadRequest, "Validation error", err)
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
			h.sendError(w, http.StatusBadRequest, "Invalid end_date format: must be mm-yyyy", err)
			return
		}
		sub.EndDate = parsedDate
	}

	logger.Info("[HANDLER] Using service...")
	if err := h.service.UpdateSubscription(r.Context(), id, sub); err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			h.sendError(w, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Error while updating subscription", err)
		return
	}

	logger.Info("[HANDLER] Creating subscription response...")
	response := dto.ToUpdateSubscriptionResponse(sub)
	h.sendSuccess(w, http.StatusOK, response)
	logger.Info("[HANDLER] Subscription updated successfully")
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID" minimum(1)
// @Success 204 "No Content - Subscription deleted successfully"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start deleting subscription...")

	id, err := getIdParameter(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	logger.Info("[HANDLER] Using service...")
	err = h.service.DeleteSubscription(r.Context(), id)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			h.sendError(w, http.StatusNotFound, "Subscription not found", err)
			return
		}
		h.sendError(w, http.StatusInternalServerError, "Error while deleting subscription", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	logger.Info("[HANDLER] Subscription deleted successfully")
}

// CollectStats godoc
// @Summary Collect subscription statistics
// @Description Collect statistics for subscriptions within a date range
// @Tags statistics
// @Accept json
// @Produce json
// @Param request body dto.SumSubscriptionRequest true "Statistics collection request"
// @Success 200 {object} dto.SumSubscriptionResponse "Statistics collected successfully"
// @Router /subscriptions/stats [post]
func (h *SubscriptionHandler) CollectStats(w http.ResponseWriter, r *http.Request) {
	logger.Info("[HANDLER] Start collecting stats...")
	defer r.Body.Close()

	logger.Info("[HANDLER] Decoding request...")
	var req dto.SumSubscriptionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Error while decoding request body", err)
		return
	}

	logger.Info("[HANDLER] Validating request...")
	if err := req.Validate(); err != nil {
		if validationErr, ok := err.(*dto.ValidationError); ok {
			h.sendValidationError(w, validationErr)
			return
		}
		h.sendError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	var filter = &model.SubscriptionFilter{}

	if req.ServiceName != nil {
		filter.ServiceName = req.ServiceName
	}
	if req.UserID != nil {
		userUuid, err := uuid.Parse(*req.UserID)
		if err != nil {
			h.sendError(w, http.StatusBadRequest, "Invalid user_id format", err)
			return
		}
		filter.UserID = &userUuid
	}

	parsedStartPeriod, err := dto.ParseDate(&req.StartPeriod)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid start_period format: must be mm-yyyy", err)
		return
	}
	filter.StartPeriod = *parsedStartPeriod

	parsedEndPeriod, err := dto.ParseDate(&req.EndPeriod)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid end_period format: must be mm-yyyy", err)
		return
	}
	filter.EndPeriod = *parsedEndPeriod

	logger.Info("[HANDLER] Using service...")
	stats, err := h.service.CollectStats(r.Context(), filter)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Error while collecting stats", err)
		return
	}

	logger.Info("[HANDLER] Creating stats response...")
	response := dto.ToSumSubscriptionResponse(stats)
	h.sendSuccess(w, http.StatusOK, response)
	logger.Info("[HANDLER] Subscriptions stats collected")
}
