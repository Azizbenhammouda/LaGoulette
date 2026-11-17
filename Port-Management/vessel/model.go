package vessel

import (
	"time"

	"github.com/google/uuid"
)

type Vessel struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name       string     `json:"name" gorm:"not null"`
	IMO        string     `json:"imo" gorm:"unique;not null"`
	Status     bool       `json:"status"` // mazal fl port wala
	Flag       string     `json:"flag"`
	VesselType string     `json:"vessel_type"`
	Arrival    *time.Time `json:"arrival"`
	Departure  *time.Time `json:"departure"`
}
