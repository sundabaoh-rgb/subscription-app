package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"subServ/internal/domain"
	"time"

	"github.com/google/uuid"
)

type createSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid user_id format")
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATE", "start_date must be in MM-YYYY format")
		return
	}

	sub := domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_DATE", "end_date must be in MM-YYYY format")
			return
		}
		sub.EndDate = &endDate
	}

	result, err := h.service.Create(r.Context(), sub)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid id format")
		return
	}

	result, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "subscription not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid id format")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "subscription not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid id format")
		return
	}

	var req createSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "invalid request body")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid user_id format")
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATE", "start_date must be in MM-YYYY format")
		return
	}

	sub := domain.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_DATE", "end_date must be in MM-YYYY format")
			return
		}
		sub.EndDate = &endDate
	}

	err = h.service.Update(r.Context(), sub)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "subscription not found")
			return
		}
		if errors.Is(err, domain.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	// 1. читаем query параметры
	filter := domain.SubscriptionFilter{
		Page:  1,
		Limit: 20,
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid user_id format")
			return
		}
		filter.UserID = &userID
	}

	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			writeError(w, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		filter.Page = page
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 100 {
			writeError(w, http.StatusBadRequest, "INVALID_LIMIT", "limit must be between 1 and 100")
			return
		}
		filter.Limit = limit
	}

	result, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	filter := domain.TotalCostFilter{}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_UUID", "invalid user_id format")
			return
		}
		filter.UserID = &userID
	}

	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		from, err := time.Parse("01-2006", fromStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_DATE", "from must be in MM-YYYY format")
			return
		}
		filter.From = &from
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		to, err := time.Parse("01-2006", toStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_DATE", "to must be in MM-YYYY format")
			return
		}
		filter.To = &to
	}

	total, err := h.service.TotalCost(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total_cost": total,
		"currency":   "RUB",
		"period": map[string]any{
			"from": r.URL.Query().Get("from"),
			"to":   r.URL.Query().Get("to"),
		},
	})
}
