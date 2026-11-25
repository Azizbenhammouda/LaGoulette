package container

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusExpected  = "EXPECTED"
	StatusUnloaded  = "UNLOADED"
	StatusStored    = "STORED"
	StatusCollected = "COLLECTED"
)

var (
	ErrContainerExists   = errors.New("container already exists")
	ErrContainerNotFound = errors.New("container not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrCannotDelete      = errors.New("only expected containers can be deleted")
)

var allowedTransitions = map[string]string{
	StatusExpected: StatusUnloaded,
	StatusUnloaded: StatusStored,
	StatusStored:   StatusCollected,
}

type RegisterContainerInput struct {
	CargoType       string    `json:"cargo_type"`
	VesselID        uuid.UUID `json:"vessel_id"`
	ContainerNumber string    `json:"container_number"`
}
type UpdateContainerInput struct {
	VesselID  *uuid.UUID `json:"vessel_id"`
	CargoType *string    `json:"cargo_type"`
}
type ChangeStatusInput struct {
	Status string `json:"status"`
}
type ContainerService interface {
	RegisterContainer(RegisterContainerInput) (*Container, error)
	GetContainer(uuid.UUID) (*Container, error)
	ListContainers() ([]Container, error)
	SearchContainers(ContainerFilter) ([]Container, error)
	UpdateContainer(uuid.UUID, UpdateContainerInput) error
	ChangeStatus(uuid.UUID, string) error
	DeleteContainer(uuid.UUID) error
}
type containerService struct {
	repo ContainerRepository
}

func NewContainerService(repo ContainerRepository) ContainerService {
	return containerService{
		repo: repo,
	}
}

func validCargoType(t string) bool {
	switch t {
	case "GENERAL", "REEFER", "HAZARDOUS", "BULK":
		return true
	}
	return false
}

func (cs containerService) RegisterContainer(ri RegisterContainerInput) (*Container, error) {
	if ri.ContainerNumber == "" || ri.VesselID == uuid.Nil || !validCargoType(ri.CargoType) {
		return nil, ErrInvalidInput
	}
	c, err := cs.repo.GetByNumber(ri.ContainerNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if c != nil {
		return nil, ErrContainerExists
	}
	con := Container{
		ContainerNumber: ri.ContainerNumber,
		VesselID:        ri.VesselID,
		CargoType:       ri.CargoType,
		Status:          StatusExpected,
	}
	err = cs.repo.Create(&con)
	if err != nil {
		return nil, err
	}
	return &con, nil
}
func (cs containerService) GetContainer(id uuid.UUID) (*Container, error) {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContainerNotFound
		} else {
			return nil, err
		}
	}
	return c, nil
}
func (cs containerService) ListContainers() ([]Container, error) {
	cons, err := cs.repo.List()
	return cons, err
}
func (cs containerService) SearchContainers(filter ContainerFilter) ([]Container, error) {
	cons, err := cs.repo.ListContainers(filter)
	return cons, err
}
func (cs containerService) UpdateContainer(id uuid.UUID, up UpdateContainerInput) error {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContainerNotFound
		} else {
			return err
		}
	}
	if up.CargoType != nil {
		if !validCargoType(*up.CargoType) {
			return ErrInvalidInput
		}
		c.CargoType = *up.CargoType
	}
	if up.VesselID != nil {
		if *up.VesselID == uuid.Nil {
			return ErrInvalidInput
		}
		c.VesselID = *up.VesselID
	}
	err = cs.repo.Update(c)
	if err != nil {
		return err
	}
	return nil
}
func (cs containerService) ChangeStatus(id uuid.UUID, newStatus string) error {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContainerNotFound
		} else {
			return err
		}
	}
	next, ok := allowedTransitions[c.Status]
	if !ok || next != newStatus {
		return ErrInvalidTransition
	}
	return cs.repo.UpdateStatus(id, newStatus)
}
func (cs containerService) DeleteContainer(id uuid.UUID) error {
	c, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContainerNotFound
		} else {
			return err
		}
	}
	if c.Status != StatusExpected {
		return ErrCannotDelete
	}
	err = cs.repo.Delete(id)
	return err
}
