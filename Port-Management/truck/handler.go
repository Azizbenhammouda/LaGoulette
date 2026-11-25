package truck

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type appointmentHandler struct {
	service AppointmentService
}

func NewAppointmentHandler(service AppointmentService) appointmentHandler {
	return appointmentHandler{
		service: service,
	}
}
func (ah appointmentHandler) RegisterAppointment(w http.ResponseWriter, req *http.Request) {
	var input RegisterAppointmentInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	a, err := ah.service.RegisterAppointment(input)
	if err != nil {
		if errors.Is(err, ErrInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrContainerNotFound) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, ErrContainerNotStored) ||
			errors.Is(err, ErrActiveAppointmentExists) ||
			errors.Is(err, ErrSlotFull) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(a)
}
func (ah appointmentHandler) ListAppointments(w http.ResponseWriter, req *http.Request) {
	as, err := ah.service.ListAppointments()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(as)
}
func (ah appointmentHandler) GetAppointment(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	a, err := ah.service.GetAppointment(id)
	if err != nil {
		if errors.Is(err, ErrAppointmentNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a)
}
func (ah appointmentHandler) Delete(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ah.service.DeleteAppointment(id)
	if err != nil {
		if errors.Is(err, ErrAppointmentNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, ErrCannotDelete) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}
func (ah appointmentHandler) UpdateAppointment(w http.ResponseWriter, req *http.Request) {
	var input UpdateAppointmentInput
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ah.service.UpdateAppointment(id, input)
	if err != nil {
		if errors.Is(err, ErrAppointmentNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, ErrInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrNotEditable) || errors.Is(err, ErrSlotFull) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (ah appointmentHandler) ChangeStatus(w http.ResponseWriter, req *http.Request) {
	var input ChangeStatusInput
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ah.service.ChangeStatus(id, input.Status)
	if err != nil {
		if errors.Is(err, ErrAppointmentNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, ErrInvalidInput) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, ErrContainerNotFound) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, ErrInvalidTransition) || errors.Is(err, ErrContainerNotStored) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (ah appointmentHandler) SearchAppointments(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	status := query.Get("status")
	idstr := query.Get("container_id")
	datestr := query.Get("date") // format: 2026-11-25
	filter := AppointmentFilter{}
	if status != "" {
		filter.Status = &status
	}
	if idstr != "" {
		containerid, err := uuid.Parse(idstr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.ContainerID = &containerid
	}
	if datestr != "" {
		date, err := time.Parse("2006-01-02", datestr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Date = &date
	}
	as, err := ah.service.SearchAppointments(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(as)
}

// GET /containers/{id}/appointments
func (ah appointmentHandler) ListByContainer(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	as, err := ah.service.SearchAppointments(AppointmentFilter{ContainerID: &id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(as)
}
