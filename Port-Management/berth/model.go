package berth

import "github.com/google/uuid"

type Berth struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name     string    `json:"name" gorm:"unique;not null"`
	Terminal string    `json:"terminal"`
	Status   string    `json:"status"`
}
