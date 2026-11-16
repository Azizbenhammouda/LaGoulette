package vessel

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VesselRepository interface {
	Create(*Vessel) error
	GetByID(uuid.UUID) (*Vessel, error)
	List() ([]Vessel, error)
	Update(*Vessel) error
	Delete(uuid.UUID) error
}

type vesselRepository struct {
	db *gorm.DB
}

func NewVesselRepository(db *gorm.DB) VesselRepository {
	return vesselRepository{
		db: db,
	}
}
func (vr vesselRepository) Create(v *Vessel) error {
	return vr.db.Create(v).Error
}
func (vr vesselRepository) GetByID(id uuid.UUID) (*Vessel, error) {
	var v Vessel
	result := vr.db.Where("id = ?", id).First(&v)
	if result.Error != nil {
		return nil, result.Error
	}
	return &v, nil
}
func (vr vesselRepository) List() ([]Vessel, error) {
	var vs []Vessel
	if err := vr.db.Where("status = ?", true).Find(&vs).Error; err != nil {
		return nil, err
	}
	return vs, nil
}
func (vr vesselRepository) Update(v *Vessel) error {
	return vr.db.Save(v).Error
}

func (vr vesselRepository) Delete(id uuid.UUID) error {
	result := vr.db.Where("id = ?", id).Delete(&Vessel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
