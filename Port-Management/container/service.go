package container

import (
	"errors"

	"github.com/google/uuid"
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

var (
	ErrContainerExists   = errors.New("container already exists")
	ErrContainerNotFound = errors.New("container not found")
)
