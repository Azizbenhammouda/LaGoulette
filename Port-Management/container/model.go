package container

import (
	"github.com/google/uuid"
)

type Container struct {
	ID              uuid.UUID `json:"id"  gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ContainerNumber string    `json:"container_number" gorm:"unique"`
	VesselID        uuid.UUID `json:"vessel_id"`
	CargoType       string    `json:"cargo_type"`
	Status          string    `json:"status"` // EXPECTED, UNLOADED, STORED, COLLECTED
}
