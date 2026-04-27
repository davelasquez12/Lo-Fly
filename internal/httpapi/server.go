package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"LoflyBE/internal/tracker"
	"LoflyBE/internal/tracker/repository"
	"LoflyBE/internal/tracker/service"
)

type Server struct {
	service service.ITrackerService
	repo    repository.ITrackerRepo
}

func NewServer(service service.ITrackerService, repo repository.ITrackerRepo) *Server {
	return &Server{service: service, repo: repo}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /tracked-flights", s.getTrackedFlights)
	mux.HandleFunc("POST /tracked-flights", s.createTrackedFlight)
	mux.HandleFunc("GET /tracked-flights/{id}", s.getTrackedFlight)
	mux.HandleFunc("DELETE /tracked-flights/{id}", s.deleteTrackedFlight)
	return withCORS(mux)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) getTrackedFlights(w http.ResponseWriter, r *http.Request) {
	trackedFlights, err := s.repo.ListTrackedFlights()
	respond(w, trackedFlights, err)
}

func (s *Server) createTrackedFlight(w http.ResponseWriter, r *http.Request) {
	var input tracker.CreateTrackedFlightRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "request body must be valid JSON")
		return
	}
	trackedFlight, err := s.service.CreateTrackedFlight(input)
	respondWithStatus(w, http.StatusCreated, trackedFlight, err)
}

func (s *Server) getTrackedFlight(w http.ResponseWriter, r *http.Request) {
	trackedFlight, err := s.repo.GetTrackedFlight(r.PathValue("id"))
	if err != nil {
		respondError(w, err)
		return
	}
	if trackedFlight == nil {
		respondError(w, tracker.ErrNotFound)
		return
	}
	writeJSON(w, http.StatusOK, envelope(trackedFlight))
}

func (s *Server) deleteTrackedFlight(w http.ResponseWriter, r *http.Request) {
	deleted, err := s.repo.DeleteTrackedFlight(r.PathValue("id"))
	if err != nil {
		respondError(w, err)
		return
	}
	if !deleted {
		respondError(w, tracker.ErrNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respond(w http.ResponseWriter, data interface{}, err error) {
	respondWithStatus(w, http.StatusOK, data, err)
}

func respondWithStatus(w http.ResponseWriter, status int, data interface{}, err error) {
	if err != nil {
		respondError(w, err)
		return
	}
	writeJSON(w, status, envelope(data))
}

func respondError(w http.ResponseWriter, err error) {
	status := statusForError(err)
	message := err.Error()
	if status == http.StatusInternalServerError {
		message = "internal server error"
	}
	writeError(w, status, message)
}

func statusForError(err error) int {
	switch {
	case errors.Is(err, tracker.ErrInvalidRequest):
		return http.StatusBadRequest
	case errors.Is(err, tracker.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func envelope(data interface{}) map[string]interface{} {
	return map[string]interface{}{"data": data}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{"message": message},
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
