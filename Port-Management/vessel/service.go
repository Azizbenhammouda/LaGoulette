package vessel

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInvalidInput = errors.New("invalid Input")
var ErrVesselNotFound = errors.New("vessel not found")
var ErrVesselExists = errors.New("vessel is already in the port")

type RegisterInput struct {
	IMO  string `json:"imo"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}
type VesselService interface {
	Register(RegisterInput) (*Vessel, error)
	GetVessel(uuid.UUID) (*Vessel, error)
	ListVessels() ([]Vessel, error)
	Delete(uuid.UUID) error
}
type vesselService struct {
	repo VesselRepository
}

func NewVesselService(repo VesselRepository) VesselService {
	return vesselService{
		repo: repo,
	}
}
func (vs vesselService) Register(input RegisterInput) (*Vessel, error) {
	existingVessel, err := vs.repo.GetByIMO(input.IMO)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingVessel != nil {
		return nil, ErrVesselExists
	}
	v := Vessel{
		IMO:  input.IMO,
		Name: input.Name,
		Flag: input.Flag,
	}
	err = vs.repo.Create(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
func (vs vesselService) GetVessel(id uuid.UUID) (*Vessel, error) {
	v, err := vs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVesselNotFound
		}
		return nil, err
	}
	return v, nil
}
func (vs vesselService) ListVessels() ([]Vessel, error) {
	return vs.repo.List()
}
func (vs vesselService) Delete(id uuid.UUID) error {
	err := vs.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVesselNotFound
		}
		return err
	}
	return nil
}
