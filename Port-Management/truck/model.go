package truck

import (
	"time"

	"github.com/google/uuid"
)

type TruckAppointment struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ContainerID uuid.UUID `json:"container_id"`
	TruckPlate  string    `json:"truck_plate"`
	DriverName  string    `json:"driver_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Status      string    `json:"status"` // SCHEDULED, ARRIVED, COMPLETED, CANCELLED
}
