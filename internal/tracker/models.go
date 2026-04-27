package tracker

import (
	"time"
)

type TrackedFlight struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	User           User           `json:"user"`
	SearchCriteria SearchCriteria `json:"searchCriteria"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type User struct {
	Email string `json:"email"`
}

type SearchCriteria struct {
	Slices         []SliceRequest     `json:"slices"`
	Passengers     []PassengerRequest `json:"passengers"`
	CabinClass     string             `json:"cabin_class,omitempty"`
	MaxConnections int                `json:"max_connections"`
}

type SliceRequest struct {
	Origin        string `json:"origin"`
	Destination   string `json:"destination"`
	DepartureDate string `json:"departure_date"`
}

type PassengerRequest struct {
	Type string `json:"type,omitempty"`
	Age  *int   `json:"age,omitempty"`
}

type CreateTrackedFlightRequest struct {
	Name           string             `json:"name"`
	User           User               `json:"user"`
	Slices         []SliceRequest     `json:"slices"`
	Passengers     []PassengerRequest `json:"passengers"`
	CabinClass     string             `json:"cabin_class"`
	MaxConnections *int               `json:"max_connections"`
}
