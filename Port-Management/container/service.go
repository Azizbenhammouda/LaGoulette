package container

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrContainerExists   = errors.New("container already exists")
	ErrContainerNotFound = errors.New("container not found")
)

type RegisterContainerInput struct {
	CargoType       string    `json:"cargo_type"`
	VesselID        uuid.UUID `json:"vessel_id"`
	ContainerNumber string    `json:"container_number"`
	Status          string    `json:"status"`
}
type UpdateContainerInput struct {
	VesselID  *uuid.UUID `json:"vessel_id"`
	CargoType *string    `json:"cargo_type"`
	Status    *string    `json:"status"`
}
type ContainerService interface {
	RegisterContainer(RegisterContainerInput) (*Container, error)
	GetContainer(uuid.UUID) (*Container, error)
	ListContainers() ([]Container, error)
	SearchContainers(ContainerFilter) ([]Container, error)
	UpdateContainer(uuid.UUID, UpdateContainerInput) error
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

func (cs containerService) RegisterContainer(ri RegisterContainerInput) (*Container, error) {
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
		Status:          ri.Status,
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
	c.CargoType = *up.CargoType
	c.VesselID = *up.VesselID
	c.Status = *up.Status
	err = cs.repo.Update(c)
	if err != nil {
		return err
	}
	return nil
}
func (cs containerService) DeleteContainer(id uuid.UUID) error {
	_, err := cs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContainerNotFound
		} else {
			return err
		}
	}
	err = cs.repo.Delete(id)
	return err
}
