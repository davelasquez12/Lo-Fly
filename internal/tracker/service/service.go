package service

import (
	"fmt"
	"time"

	"LoflyBE/internal/tracker"
	"LoflyBE/internal/tracker/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock_tracker_service.go -package=mocks

type ITrackerService interface {
	CreateTrackedFlight(input tracker.CreateTrackedFlightRequest) (*tracker.TrackedFlight, error)
}

type Service struct {
	repo repository.ITrackerRepo
}

func NewService(repo repository.ITrackerRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrackedFlight(input tracker.CreateTrackedFlightRequest) (*tracker.TrackedFlight, error) {
	if err := validateCreateRequest(input); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	maxConnections := 1
	if input.MaxConnections != nil {
		maxConnections = *input.MaxConnections
	}
	trackedFlight := tracker.TrackedFlight{
		ID:   newID("trk"),
		Name: defaultName(input.Slices),
		User: input.User,
		SearchCriteria: tracker.SearchCriteria{
			Slices:         input.Slices,
			Passengers:     input.Passengers,
			CabinClass:     input.CabinClass,
			MaxConnections: maxConnections,
		},
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if input.Name != "" {
		trackedFlight.Name = input.Name
	}

	if err := s.repo.CreateTrackedFlight(trackedFlight); err != nil {
		return nil, err
	}
	return &trackedFlight, nil
}

func validateCreateRequest(input tracker.CreateTrackedFlightRequest) error {
	if len(input.Slices) == 0 {
		return fmt.Errorf("%w: at least one slice is required", tracker.ErrInvalidRequest)
	}
	if len(input.Passengers) == 0 {
		return fmt.Errorf("%w: at least one passenger is required", tracker.ErrInvalidRequest)
	}
	for _, slice := range input.Slices {
		if slice.Origin == "" || slice.Destination == "" || slice.DepartureDate == "" {
			return fmt.Errorf("%w: each slice requires origin, destination, and departure_date", tracker.ErrInvalidRequest)
		}
	}
	return nil
}

func defaultName(slices []tracker.SliceRequest) string {
	if len(slices) == 0 {
		return "Tracked flight"
	}
	first := slices[0]
	return first.Origin + "-" + first.Destination
}
