package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"LoflyBE/internal/tracker"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock_tracker_repo.go -package=mocks

type ITrackerRepo interface {
	ListTrackedFlights() ([]tracker.TrackedFlight, error)
	GetTrackedFlight(id string) (*tracker.TrackedFlight, error)
	CreateTrackedFlight(trackedFlight tracker.TrackedFlight) error
	UpdateTrackedFlight(trackedFlight tracker.TrackedFlight) error
	DeleteTrackedFlight(id string) (bool, error)
}

type TrackerRepo struct {
	path string
	mu   sync.Mutex
}

type store struct {
	TrackedFlights []tracker.TrackedFlight `json:"trackedFlights"`
}

func NewTrackerRepo(path string) *TrackerRepo {
	return &TrackerRepo{path: path}
}

func (r *TrackerRepo) ListTrackedFlights() ([]tracker.TrackedFlight, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.read()
	if err != nil {
		return nil, err
	}
	return append([]tracker.TrackedFlight(nil), data.TrackedFlights...), nil
}

func (r *TrackerRepo) GetTrackedFlight(id string) (*tracker.TrackedFlight, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.read()
	if err != nil {
		return nil, err
	}
	for _, trackedFlight := range data.TrackedFlights {
		if trackedFlight.ID == id {
			copy := trackedFlight
			return &copy, nil
		}
	}
	return nil, nil
}

func (r *TrackerRepo) CreateTrackedFlight(trackedFlight tracker.TrackedFlight) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.read()
	if err != nil {
		return err
	}
	data.TrackedFlights = append(data.TrackedFlights, trackedFlight)
	return r.write(data)
}

func (r *TrackerRepo) UpdateTrackedFlight(trackedFlight tracker.TrackedFlight) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.read()
	if err != nil {
		return err
	}
	for i := range data.TrackedFlights {
		if data.TrackedFlights[i].ID == trackedFlight.ID {
			data.TrackedFlights[i] = trackedFlight
			return r.write(data)
		}
	}
	return errors.New("tracked flight not found")
}

func (r *TrackerRepo) DeleteTrackedFlight(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, err := r.read()
	if err != nil {
		return false, err
	}
	next := data.TrackedFlights[:0]
	deleted := false
	for _, trackedFlight := range data.TrackedFlights {
		if trackedFlight.ID == id {
			deleted = true
			continue
		}
		next = append(next, trackedFlight)
	}
	data.TrackedFlights = next
	return deleted, r.write(data)
}

func (r *TrackerRepo) read() (store, error) {
	raw, err := os.ReadFile(r.path)
	if os.IsNotExist(err) {
		return store{}, nil
	}
	if err != nil {
		return store{}, err
	}
	var data store
	if err := json.Unmarshal(raw, &data); err != nil {
		return store{}, err
	}
	return data, nil
}

func (r *TrackerRepo) write(data store) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(r.path, raw, 0644)
}
