package service_test

import (
	"errors"
	"testing"

	"LoflyBE/internal/tracker"
	"LoflyBE/internal/tracker/repository/mocks"
	"LoflyBE/internal/tracker/service"

	"go.uber.org/mock/gomock"
)

func TestCreateTrackedFlightDefaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockITrackerRepo(ctrl)
	svc := service.NewService(repo)

	var stored tracker.TrackedFlight
	repo.EXPECT().
		CreateTrackedFlight(gomock.Any()).
		DoAndReturn(func(trackedFlight tracker.TrackedFlight) error {
			stored = trackedFlight
			return nil
		})

	trackedFlight, err := svc.CreateTrackedFlight(tracker.CreateTrackedFlightRequest{
		Slices: []tracker.SliceRequest{{
			Origin:        "ORD",
			Destination:   "LHR",
			DepartureDate: "2026-09-10",
		}},
		Passengers: []tracker.PassengerRequest{{Type: "adult"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if stored.ID != trackedFlight.ID {
		t.Fatalf("expected stored tracked flight ID %s, got %s", trackedFlight.ID, stored.ID)
	}
	if trackedFlight.Name != "ORD-LHR" {
		t.Fatalf("expected default name ORD-LHR, got %s", trackedFlight.Name)
	}
	if trackedFlight.SearchCriteria.MaxConnections != 1 {
		t.Fatalf("expected max_connections default 1, got %d", trackedFlight.SearchCriteria.MaxConnections)
	}
	if trackedFlight.Status != "active" {
		t.Fatalf("expected active status, got %s", trackedFlight.Status)
	}
}

func TestCreateTrackedFlightValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockITrackerRepo(ctrl)
	svc := service.NewService(repo)

	_, err := svc.CreateTrackedFlight(tracker.CreateTrackedFlightRequest{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !errors.Is(err, tracker.ErrInvalidRequest) {
		t.Fatalf("expected invalid request error, got %v", err)
	}
}

func TestCreateTrackedFlightReturnsRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockITrackerRepo(ctrl)
	svc := service.NewService(repo)

	repoErr := errors.New("repo failed")
	repo.EXPECT().
		CreateTrackedFlight(gomock.Any()).
		Return(repoErr)

	_, err := svc.CreateTrackedFlight(tracker.CreateTrackedFlightRequest{
		Slices: []tracker.SliceRequest{{
			Origin:        "ORD",
			Destination:   "LHR",
			DepartureDate: "2026-09-10",
		}},
		Passengers: []tracker.PassengerRequest{{Type: "adult"}},
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
