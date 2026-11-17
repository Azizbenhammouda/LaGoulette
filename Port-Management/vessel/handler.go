package vessel

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type vesselHandler struct {
	service VesselService
}

func NewVesselHandler(service VesselService) vesselHandler {
	return vesselHandler{
		service: service,
	}
}

func (vh vesselHandler) Register(w http.ResponseWriter, req *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	vessel, err := vh.service.Register(input)
	if err != nil {
		if errors.Is(err, ErrVesselExists) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vessel)
}
func (vh vesselHandler) List(w http.ResponseWriter, req *http.Request) {
	vessels, err := vh.service.ListVessels()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vessels)
}
func (vh vesselHandler) GetVessel(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	vesselID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	v, err := vh.service.GetVessel(vesselID)
	if err != nil {
		if errors.Is(err, ErrVesselNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}
func (vh vesselHandler) Update(w http.ResponseWriter, req *http.Request) {
	var up RequestUpdate
	idStr := req.PathValue("id")
	vesselID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(req.Body).Decode(&up); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = vh.service.Update(vesselID, up)
	if err != nil {
		if errors.Is(err, ErrVesselNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (vh vesselHandler) Delete(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	vesselID, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = vh.service.Delete(vesselID)
	if err != nil {
		if errors.Is(err, ErrVesselNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
