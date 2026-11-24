package container

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type containerHandler struct {
	service ContainerService
}

func NewContainerHandler(service ContainerService) containerHandler {
	return containerHandler{
		service: service,
	}
}
func (ch containerHandler) RegisterContainer(w http.ResponseWriter, req *http.Request) {
	var input RegisterContainerInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c, err := ch.service.RegisterContainer(input)
	if err != nil {
		if errors.Is(err, ErrContainerExists) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}
func (ch containerHandler) ListContainers(w http.ResponseWriter, req *http.Request) {
	cs, err := ch.service.ListContainers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cs)
}
func (ch containerHandler) GetContainer(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c, err := ch.service.GetContainer(id)
	if err != nil {
		if errors.Is(err, ErrContainerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}
func (ch containerHandler) Delete(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = ch.service.DeleteContainer(id)
	if err != nil {
		if errors.Is(err, ErrContainerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (ch containerHandler) UpdateContainer(w http.ResponseWriter, req *http.Request) {
	var input UpdateContainerInput
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
	err = ch.service.UpdateContainer(id, input)
	if err != nil {
		if errors.Is(err, ErrContainerNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (ch containerHandler) SearchContainers(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	status := query.Get("status")
	cargotype := query.Get("cargo_type")
	idstr := query.Get("vessel_id")
	vesselid, err := uuid.Parse(idstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := ContainerFilter{}
	if status != "" {
		filter.Status = &status
	}
	if cargotype != "" {
		filter.CargoType = &cargotype
	}
	if vesselid != uuid.Nil {
		filter.VesselID = &vesselid
	}
	cons, err := ch.service.SearchContainers(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cons)
}
