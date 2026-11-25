package truck

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"PortManagement/container"
)

const (
	StatusScheduled = "SCHEDULED"
	StatusArrived   = "ARRIVED"
	StatusCompleted = "COMPLETED"
	StatusCancelled = "CANCELLED"

	MaxTrucksPerSlot = 10
)

var (
	ErrAppointmentNotFound     = errors.New("appointment not found")
	ErrInvalidInput            = errors.New("invalid input")
	ErrContainerNotFound       = errors.New("container not found")
	ErrContainerNotStored      = errors.New("container is not in STORED status")
	ErrActiveAppointmentExists = errors.New("container already has an active appointment")
	ErrSlotFull                = errors.New("time slot is full")
	ErrInvalidTransition       = errors.New("invalid status transition")
	ErrNotEditable             = errors.New("only scheduled appointments can be edited")
	ErrCannotDelete            = errors.New("only scheduled or cancelled appointments can be deleted")
)

var allowedTransitions = map[string][]string{
	StatusScheduled: {StatusArrived, StatusCancelled},
	StatusArrived:   {StatusCompleted, StatusCancelled},
}

type RegisterAppointmentInput struct {
	ContainerID uuid.UUID `json:"container_id"`
	TruckPlate  string    `json:"truck_plate"`
	DriverName  string    `json:"driver_name"`
	ScheduledAt time.Time `json:"scheduled_at"`
}
type UpdateAppointmentInput struct {
	TruckPlate  *string    `json:"truck_plate"`
	DriverName  *string    `json:"driver_name"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}
type ChangeStatusInput struct {
	Status string `json:"status"`
}

type ContainerProvider interface {
	GetContainer(uuid.UUID) (*container.Container, error)
	ChangeStatus(uuid.UUID, string) error
}
type AppointmentService interface {
	RegisterAppointment(RegisterAppointmentInput) (*TruckAppointment, error)
	GetAppointment(uuid.UUID) (*TruckAppointment, error)
	ListAppointments() ([]TruckAppointment, error)
	SearchAppointments(AppointmentFilter) ([]TruckAppointment, error)
	UpdateAppointment(uuid.UUID, UpdateAppointmentInput) error
	ChangeStatus(uuid.UUID, string) error
	CancelAppointment(uuid.UUID) error
	DeleteAppointment(uuid.UUID) error
}
type appointmentService struct {
	repo       AppointmentRepository
	containers ContainerProvider
}

func NewAppointmentService(repo AppointmentRepository, containers ContainerProvider) AppointmentService {
	return appointmentService{
		repo:       repo,
		containers: containers,
	}
}

func validStatus(s string) bool {
	switch s {
	case StatusScheduled, StatusArrived, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}
func canTransition(from, to string) bool {
	for _, next := range allowedTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func (as appointmentService) RegisterAppointment(ri RegisterAppointmentInput) (*TruckAppointment, error) {
	if ri.ContainerID == uuid.Nil || ri.TruckPlate == "" || !ri.ScheduledAt.After(time.Now()) {
		return nil, ErrInvalidInput
	}
	slot := ri.ScheduledAt.UTC().Truncate(time.Hour)

	c, err := as.containers.GetContainer(ri.ContainerID)
	if err != nil {
		if errors.Is(err, container.ErrContainerNotFound) {
			return nil, ErrContainerNotFound
		} else {
			return nil, err
		}
	}
	if c.Status != container.StatusStored {
		return nil, ErrContainerNotStored
	}

	exists, err := as.repo.ExistsActiveForContainer(ri.ContainerID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrActiveAppointmentExists
	}

	count, err := as.repo.CountByTimeSlot(slot)
	if err != nil {
		return nil, err
	}
	if count >= MaxTrucksPerSlot {
		return nil, ErrSlotFull
	}

	a := TruckAppointment{
		ContainerID: ri.ContainerID,
		TruckPlate:  ri.TruckPlate,
		DriverName:  ri.DriverName,
		ScheduledAt: slot,
		Status:      StatusScheduled,
	}
	err = as.repo.Create(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
func (as appointmentService) GetAppointment(id uuid.UUID) (*TruckAppointment, error) {
	a, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppointmentNotFound
		} else {
			return nil, err
		}
	}
	return a, nil
}
func (as appointmentService) ListAppointments() ([]TruckAppointment, error) {
	list, err := as.repo.List()
	return list, err
}
func (as appointmentService) SearchAppointments(filter AppointmentFilter) ([]TruckAppointment, error) {
	list, err := as.repo.ListAppointments(filter)
	return list, err
}
func (as appointmentService) UpdateAppointment(id uuid.UUID, up UpdateAppointmentInput) error {
	a, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAppointmentNotFound
		} else {
			return err
		}
	}
	if a.Status != StatusScheduled {
		return ErrNotEditable
	}
	if up.TruckPlate != nil {
		if *up.TruckPlate == "" {
			return ErrInvalidInput
		}
		a.TruckPlate = *up.TruckPlate
	}
	if up.DriverName != nil {
		a.DriverName = *up.DriverName
	}
	if up.ScheduledAt != nil {
		if !up.ScheduledAt.After(time.Now()) {
			return ErrInvalidInput
		}
		slot := up.ScheduledAt.UTC().Truncate(time.Hour)
		if !slot.Equal(a.ScheduledAt.UTC().Truncate(time.Hour)) {
			count, err := as.repo.CountByTimeSlot(slot)
			if err != nil {
				return err
			}
			if count >= MaxTrucksPerSlot {
				return ErrSlotFull
			}
		}
		a.ScheduledAt = slot
	}
	err = as.repo.Update(a)
	if err != nil {
		return err
	}
	return nil
}
func (as appointmentService) ChangeStatus(id uuid.UUID, newStatus string) error {
	if !validStatus(newStatus) {
		return ErrInvalidInput
	}
	a, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAppointmentNotFound
		} else {
			return err
		}
	}
	if !canTransition(a.Status, newStatus) {
		return ErrInvalidTransition
	}
	if newStatus == StatusCompleted {
		err = as.containers.ChangeStatus(a.ContainerID, container.StatusCollected)
		if err != nil {
			if errors.Is(err, container.ErrContainerNotFound) {
				return ErrContainerNotFound
			} else if errors.Is(err, container.ErrInvalidTransition) {
				return ErrContainerNotStored
			} else {
				return err
			}
		}
	}
	return as.repo.UpdateStatus(id, newStatus)
}
func (as appointmentService) CancelAppointment(id uuid.UUID) error {
	return as.ChangeStatus(id, StatusCancelled)
}
func (as appointmentService) DeleteAppointment(id uuid.UUID) error {
	a, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAppointmentNotFound
		} else {
			return err
		}
	}
	if a.Status != StatusScheduled && a.Status != StatusCancelled {
		return ErrCannotDelete
	}
	err = as.repo.Delete(id)
	return err
}
