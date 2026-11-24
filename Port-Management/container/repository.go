package container

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ContainerFilter struct {
	Status    *string
	VesselID  *uuid.UUID
	CargoType *string
}
type ContainerRepository interface {
	Create(*Container) error
	GetByID(id uuid.UUID) (*Container, error)
	GetByNumber(number string) (*Container, error)
	List() ([]Container, error)
	ListContainers(filter ContainerFilter) ([]Container, error)
	Update(*Container) error
	UpdateStatus(id, status string) error
	Delete(id uuid.UUID) error
}
type containerRepository struct {
	db *gorm.DB
}

func NewContainerRepository(db *gorm.DB) ContainerRepository {
	return containerRepository{
		db: db,
	}
}
func (cr containerRepository) Create(c *Container) error {
	return cr.db.Create(c).Error
}
func (cr containerRepository) GetByID(id uuid.UUID) (*Container, error) {
	var c Container
	result := cr.db.Where("id = ?", id).First(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}
func (cr containerRepository) GetByNumber(num string) (*Container, error) {
	var c Container
	result := cr.db.Where("container_number = ?", num).First(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return &c, nil
}

func (cr containerRepository) Update(c *Container) error {
	return cr.db.Save(c).Error
}
func (cr containerRepository) UpdateStatus(id, status string) error {
	result := cr.db.Model(&Container{}).
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
func (cr containerRepository) Delete(id uuid.UUID) error {
	result := cr.db.Where("id = ?", id).Delete(&Container{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (cr containerRepository) List() ([]Container, error) {
	var cs []Container
	res := cr.db.Find(&cs)
	if res.Error != nil {
		return nil, res.Error
	}
	return cs, nil
}
func (cr containerRepository) ListContainers(filter ContainerFilter) ([]Container, error) {
	var cs []Container
	query := cr.db.Model(&Container{})
	if filter.Status != nil {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CargoType != nil {
		query = query.Where("cargo_type = ?", filter.CargoType)
	}
	if filter.VesselID != nil {
		query = query.Where("vessel_id = ?", filter.VesselID)
	}
	result := query.Find(&cs)
	if result.Error != nil {
		return nil, result.Error
	}
	return cs, nil
}
