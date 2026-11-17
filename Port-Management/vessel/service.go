package vessel

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrInvalidInput = errors.New("invalid Input")
var ErrVesselNotFound = errors.New("vessel not found")
var ErrVesselExists = errors.New("vessel is already in the port")

type RegisterInput struct {
	IMO        string `json:"imo"`
	Name       string `json:"name"`
	Flag       string `json:"flag"`
	VesselType string `json:"vessel_type"`
}
type RequestUpdate struct {
	Name       *string    `json:"name,omitempty"`
	Status     *bool      `json:"status,omitempty"`
	Flag       *string    `json:"flag,omitempty"`
	VesselType *string    `json:"vessel_type,omitempty"`
	Arrival    *time.Time `json:"arrival,omitempty"`
	Departure  *time.Time `json:"departure,omitempty"`
}
type VesselService interface {
	Register(RegisterInput) (*Vessel, error)
	GetVessel(uuid.UUID) (*Vessel, error)
	ListVessels() ([]Vessel, error)
	Delete(uuid.UUID) error
	Update(uuid.UUID, RequestUpdate) error
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
	arriv := time.Now()
	v := Vessel{
		IMO:        input.IMO,
		Name:       input.Name,
		Flag:       input.Flag,
		VesselType: input.VesselType,
		Status:     true,
		Arrival:    &arriv,
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
func (vs vesselService) Update(id uuid.UUID, req RequestUpdate) error {
	v, err := vs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVesselNotFound
		}
		return err
	}
	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.Flag != nil {
		v.Flag = *req.Flag
	}
	if req.Status != nil {
		v.Status = *req.Status
	}
	if req.VesselType != nil {
		v.VesselType = *req.VesselType
	}
	if req.Departure != nil {
		v.Departure = req.Departure
	}
	if req.Arrival != nil {
		v.Arrival = req.Arrival
	}
	if err := vs.repo.Update(v); err != nil {
		return err
	}
	return nil
}
