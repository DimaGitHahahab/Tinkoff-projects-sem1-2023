package handler

import (
	"encoding/json"
	"errors"
	"homework/internal/model"
	"homework/internal/service"
	"net/http"
)

type Handler struct {
	Service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	d := model.Device{}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.ErrResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Service.CreateDevice(d); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	num := r.URL.Query().Get("num")
	d, err := h.Service.GetDevice(num)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response, err := json.Marshal(d)
	if err != nil {
		h.ErrResponse(w, "JSON can't be marshaled", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *Handler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	num := r.URL.Query().Get("num")
	if err := h.Service.DeleteDevice(num); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	d := model.Device{}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.ErrResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateDevice(d); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	var httpStatus int
	var message string

	switch {
	case errors.Is(err, service.ErrDeviceAlreadyExists):
		httpStatus = http.StatusConflict
		message = "Device already exists"
	case errors.Is(err, service.ErrDeviceDoesNotExist):
		httpStatus = http.StatusNotFound
		message = "Device doesn't exist"
	case errors.Is(err, service.ErrInvalidModel):
		fallthrough
	case errors.Is(err, service.ErrInvalidSerialNumber):
		fallthrough
	case errors.Is(err, service.ErrInvalidIPAddress):
		httpStatus = http.StatusBadRequest
		message = err.Error()
	default:
		httpStatus = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.ErrResponse(w, message, httpStatus)
}

func (h *Handler) ErrResponse(w http.ResponseWriter, message string, errStatus int) {
	response, _ := json.Marshal(map[string]string{"message": message})
	w.WriteHeader(errStatus)
	_, _ = w.Write(response)
}
