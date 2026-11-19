package berth

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrBerthNotAvailable = errors.New("berth is not available")
	ErrBerthNotFound     = errors.New("berth does not exist")
)

type BerthService interface {
	List() ([]Berth, error)
	ListAvailable() ([]Berth, error)
	ReserveBerth(string) (*Berth, error)
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

func (bs berthService) ReserveBerth(name string) (*Berth, error) {
	_, err := bs.repo.GetByName(name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBerthNotFound
		}
		return nil, err
	}
	b, reserved, err := bs.repo.Reserve(name)
	if err != nil {
		return nil, err
	}
	if !reserved {
		return nil, ErrBerthNotAvailable
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
		return nil, ErrBerthNotAvailable
	}
	return b, nil
}
