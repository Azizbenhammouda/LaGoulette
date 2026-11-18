package berth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BerthRepository interface {
	Create(*Berth) error
	Update(*Berth) error
	Delete(uuid.UUID) error
	List() ([]Berth, error)
	GetByID(uuid.UUID) (*Berth, error)
	GetByName(string) (*Berth, error)
	ListByStatus(string) ([]Berth, error)
}

type berthRepository struct {
	db *gorm.DB
}

func NewBerthRepository(db *gorm.DB) BerthRepository {
	return berthRepository{
		db: db,
	}
}

func (br berthRepository) Create(b *Berth) error {
	return br.db.Create(b).Error
}

func (br berthRepository) Update(b *Berth) error {
	return br.db.Save(b).Error
}

func (br berthRepository) List() ([]Berth, error) {
	var berths []Berth
	if err := br.db.Find(&berths).Error; err != nil {
		return nil, err
	}

	return berths, nil
}

func (br berthRepository) ListByStatus(stat string) ([]Berth, error) {
	var berths []Berth
	if err := br.db.Where("status = ?", stat).Find(&berths).Error; err != nil {
		return nil, err
	}
	return berths, nil
}

func (br berthRepository) GetByID(id uuid.UUID) (*Berth, error) {
	var b Berth
	if err := br.db.First(&b, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (br berthRepository) GetByName(name string) (*Berth, error) {
	var b Berth
	if err := br.db.Where("name = ?", name).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (br berthRepository) Delete(id uuid.UUID) error {
	result := br.db.Delete(&Berth{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
