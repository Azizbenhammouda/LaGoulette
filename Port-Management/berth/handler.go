package berth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type berthHandler struct {
	service BerthService
}

func NewBerthHandler(service BerthService) berthHandler {
	return berthHandler{
		service: service,
	}
}

func (bh berthHandler) List(w http.ResponseWriter, req *http.Request) {
	berths, err := bh.service.List()
	if err != nil {
		http.Error(w, "Failed to retrieve berths", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(berths)
}

func (bh berthHandler) ListAvailable(w http.ResponseWriter, req *http.Request) {
	berths, err := bh.service.ListAvailable()
	if err != nil {
		http.Error(w, "Failed to retrieve available berths", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(berths)
}

func (bh berthHandler) ReserveBerth(w http.ResponseWriter, req *http.Request) {
	var input string
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	b, err := bh.service.ReserveBerth(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrBerthNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, ErrBerthNotAvailable):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Failed to reserve berth", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func (bh berthHandler) GetReserve(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("id")
	berthID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid berth ID", http.StatusBadRequest)
		return
	}
	b, err := bh.service.GetReserve(berthID)
	if err != nil {
		switch {
		case errors.Is(err, ErrBerthNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, ErrBerthNotAvailable):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "Failed to retrieve reservation", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(b)
}
