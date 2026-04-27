package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"LoflyBE/internal/httpapi"
	"LoflyBE/internal/tracker"
	repomocks "LoflyBE/internal/tracker/repository/mocks"
	trackerservice "LoflyBE/internal/tracker/service"
	servicemocks "LoflyBE/internal/tracker/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestGetTrackedFlightsUsesRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := servicemocks.NewMockITrackerService(ctrl)
	repo := repomocks.NewMockITrackerRepo(ctrl)
	server := httpapi.NewServer(service, repo)

	repo.EXPECT().
		ListTrackedFlights().
		Return([]tracker.TrackedFlight{{
			ID:     "trk_123",
			Name:   "ORD-LHR",
			Status: "active",
		}}, nil)

	req := httptest.NewRequest(http.MethodGet, "/tracked-flights", nil)
	res := httptest.NewRecorder()

	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if !strings.Contains(res.Body.String(), `"id":"trk_123"`) {
		t.Fatalf("expected tracked flight response, got %s", res.Body.String())
	}
}

func TestCreateTrackedFlightMissingRequiredFieldsReturnsBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := repomocks.NewMockITrackerRepo(ctrl)
	service := trackerservice.NewService(repo)
	server := httpapi.NewServer(service, repo)

	req := httptest.NewRequest(http.MethodPost, "/tracked-flights", strings.NewReader(`{}`))
	res := httptest.NewRecorder()

	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
	if !strings.Contains(res.Body.String(), "at least one slice is required") {
		t.Fatalf("expected validation error response, got %s", res.Body.String())
	}
}

func TestCreateTrackedFlightUsesService(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := servicemocks.NewMockITrackerService(ctrl)
	repo := repomocks.NewMockITrackerRepo(ctrl)
	server := httpapi.NewServer(service, repo)

	service.EXPECT().
		CreateTrackedFlight(gomock.Any()).
		DoAndReturn(func(input tracker.CreateTrackedFlightRequest) (*tracker.TrackedFlight, error) {
			if input.Slices[0].Origin != "ORD" {
				t.Fatalf("expected origin ORD, got %s", input.Slices[0].Origin)
			}
			return &tracker.TrackedFlight{
				ID:     "trk_123",
				Name:   "ORD-LHR",
				Status: "active",
			}, nil
		})

	body := tracker.CreateTrackedFlightRequest{
		Slices: []tracker.SliceRequest{{
			Origin:        "ORD",
			Destination:   "LHR",
			DepartureDate: "2026-09-10",
		}},
		Passengers: []tracker.PassengerRequest{{Type: "adult"}},
	}
	rawBody, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/tracked-flights", bytes.NewReader(rawBody))
	res := httptest.NewRecorder()

	server.Routes().ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}
	if !strings.Contains(res.Body.String(), `"id":"trk_123"`) {
		t.Fatalf("expected tracked flight response, got %s", res.Body.String())
	}
}
