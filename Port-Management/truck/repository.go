package truck

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentFilter struct {
	Status      *string
	ContainerID *uuid.UUID
	Date        *time.Time
}
type AppointmentRepository interface {
	Create(*TruckAppointment) error
	GetByID(id uuid.UUID) (*TruckAppointment, error)
	List() ([]TruckAppointment, error)
	ListAppointments(filter AppointmentFilter) ([]TruckAppointment, error)
	Update(*TruckAppointment) error
	UpdateStatus(id uuid.UUID, status string) error
	Delete(id uuid.UUID) error
	CountByTimeSlot(slot time.Time) (int64, error)
	ExistsActiveForContainer(containerID uuid.UUID) (bool, error)
}
type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return appointmentRepository{
		db: db,
	}
}
func (ar appointmentRepository) Create(a *TruckAppointment) error {
	return ar.db.Create(a).Error
}
func (ar appointmentRepository) GetByID(id uuid.UUID) (*TruckAppointment, error) {
	var a TruckAppointment
	result := ar.db.Where("id = ?", id).First(&a)
	if result.Error != nil {
		return nil, result.Error
	}
	return &a, nil
}
func (ar appointmentRepository) List() ([]TruckAppointment, error) {
	var as []TruckAppointment
	res := ar.db.Find(&as)
	if res.Error != nil {
		return nil, res.Error
	}
	return as, nil
}
func (ar appointmentRepository) ListAppointments(filter AppointmentFilter) ([]TruckAppointment, error) {
	var as []TruckAppointment
	query := ar.db.Model(&TruckAppointment{})
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.ContainerID != nil {
		query = query.Where("container_id = ?", *filter.ContainerID)
	}
	if filter.Date != nil {
		start := filter.Date.UTC().Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("scheduled_at >= ? AND scheduled_at < ?", start, end)
	}
	result := query.Find(&as)
	if result.Error != nil {
		return nil, result.Error
	}
	return as, nil
}
func (ar appointmentRepository) Update(a *TruckAppointment) error {
	return ar.db.Save(a).Error
}
func (ar appointmentRepository) UpdateStatus(id uuid.UUID, status string) error {
	result := ar.db.Model(&TruckAppointment{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (ar appointmentRepository) Delete(id uuid.UUID) error {
	result := ar.db.Where("id = ?", id).Delete(&TruckAppointment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ar appointmentRepository) CountByTimeSlot(slot time.Time) (int64, error) {
	var count int64
	result := ar.db.Model(&TruckAppointment{}).
		Where("scheduled_at >= ? AND scheduled_at < ?", slot, slot.Add(time.Hour)).
		Where("status <> ?", StatusCancelled).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
func (ar appointmentRepository) ExistsActiveForContainer(containerID uuid.UUID) (bool, error) {
	var count int64
	result := ar.db.Model(&TruckAppointment{}).
		Where("container_id = ?", containerID).
		Where("status IN ?", []string{StatusScheduled, StatusArrived}).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
