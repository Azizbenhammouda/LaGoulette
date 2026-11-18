package berth

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrBerthOccupied     = errors.New("berth is not reserved")
	ErrBerthNotAvailable = errors.New("berth is not available")
	ErrBerthNotFound     = errors.New("berth does not exist")
)

type ReserveInput struct {
	Name     string `json:"name"`
	Terminal string `json:"terminal"`
}

type BerthService interface {
	List() ([]Berth, error)
	ListAvailable() ([]Berth, error)
	ReserveBerth(ReserveInput) (*Berth, error)
	GetReserve(uuid.UUID) (*Berth, error)
}

type berthService struct {
	repo BerthRepository
}

func NewBerthService(repo BerthRepository) BerthService {
	return berthService{
		repo: repo,
	}
}

func (bs berthService) List() ([]Berth, error) {
	return bs.repo.List()
}

func (bs berthService) ListAvailable() ([]Berth, error) {
	return bs.repo.ListByStatus("available")
}

func (bs berthService) ReserveBerth(ri ReserveInput) (*Berth, error) {
	b, err := bs.repo.GetByName(ri.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBerthNotFound
		}
		return nil, err
	}
	if b.Status != "available" {
		return nil, ErrBerthNotAvailable
	}
	b.Status = "occupied"
	if err := bs.repo.Update(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (bs berthService) GetReserve(id uuid.UUID) (*Berth, error) {
	b, err := bs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBerthNotFound
		}
		return nil, err
	}
	if b.Status != "occupied" {
		return nil, ErrBerthOccupied
	}
	return b, nil
}
