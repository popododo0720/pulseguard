package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pulseguard/pulseguard/internal/server/store"
)

func listNotificationChannelsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channels, err := s.ListNotificationChannels()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if channels == nil {
			channels = []*store.NotificationChannel{}
		}
		writeJSON(w, http.StatusOK, channels)
	}
}

func createNotificationChannelHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ch store.NotificationChannel
		if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		ch.ID = uuid.New().String()
		if err := s.CreateNotificationChannel(&ch); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, ch)
	}
}

func updateNotificationChannelHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var ch store.NotificationChannel
		if err := json.NewDecoder(r.Body).Decode(&ch); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		ch.ID = id
		if err := s.UpdateNotificationChannel(&ch); err != nil {
			writeError(w, http.StatusNotFound, "notification channel not found")
			return
		}
		writeJSON(w, http.StatusOK, ch)
	}
}

func deleteNotificationChannelHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := s.DeleteNotificationChannel(id); err != nil {
			writeError(w, http.StatusNotFound, "notification channel not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
