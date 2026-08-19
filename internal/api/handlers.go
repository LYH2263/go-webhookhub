package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LYH2263/go-webhookhub"
)

type createBody struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Secret      string            `json:"secret"`
	Events      []string          `json:"events"`
	EventCSV    string            `json:"event_csv"`
	Enabled     *bool             `json:"enabled"`
	Headers     map[string]string `json:"headers"`
	TimeoutMS   int               `json:"timeout_ms"`
	MaxAttempts int               `json:"max_attempts"`
	Description string            `json:"description"`
}

type dispatchBody struct {
	Event string `json:"event"`
	Body  string `json:"body"`
}

func (s *Server) handleListEndpoints(w http.ResponseWriter, _ *http.Request) {
	views := s.hub.ListEndpoints()
	type pub struct {
		ID          string            `json:"id"`
		URL         string            `json:"url"`
		Events      []string          `json:"events"`
		Enabled     bool              `json:"enabled"`
		SecretLen   int               `json:"secret_len"`
		Headers     map[string]string `json:"headers,omitempty"`
		MaxAttempts int               `json:"max_attempts"`
		Description string            `json:"description,omitempty"`
	}
	out := make([]pub, 0, len(views))
	for _, v := range views {
		out = append(out, pub{
			ID:          v.ID,
			URL:         v.URL,
			Events:      v.Events,
			Enabled:     v.Enabled,
			SecretLen:   v.SecretLen,
			Headers:     v.Headers,
			MaxAttempts: v.MaxAttempts,
			Description: v.Description,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"endpoints": out})
}

func (s *Server) handleCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var in createBody
	if err := readJSON(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	events := in.Events
	if len(events) == 0 && in.EventCSV != "" {
		for _, p := range strings.Split(in.EventCSV, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				events = append(events, p)
			}
		}
	}
	enabled := true
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	var timeout time.Duration
	if in.TimeoutMS > 0 {
		timeout = time.Duration(in.TimeoutMS) * time.Millisecond
	}
	id, err := s.hub.Subscribe(webhookhub.Endpoint{
		ID:          in.ID,
		URL:         in.URL,
		Secret:      []byte(in.Secret),
		Events:      events,
		Enabled:     enabled,
		Headers:     in.Headers,
		Timeout:     timeout,
		MaxAttempts: in.MaxAttempts,
		Description: in.Description,
	})
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

func (s *Server) handleDeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.hub.Unsubscribe(id); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "id": id})
}

func (s *Server) handleEnable(w http.ResponseWriter, r *http.Request) {
	s.setEnabled(w, r, true)
}

func (s *Server) handleDisable(w http.ResponseWriter, r *http.Request) {
	s.setEnabled(w, r, false)
}

func (s *Server) setEnabled(w http.ResponseWriter, r *http.Request, enabled bool) {
	id := r.PathValue("id")
	if err := s.hub.SetEnabled(id, enabled); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "id": id, "enabled": enabled})
}

func (s *Server) handleDispatch(w http.ResponseWriter, r *http.Request) {
	var in dispatchBody
	if err := readJSON(r, &in); err != nil {
		writeErr(w, err)
		return
	}
	results, err := s.hub.DispatchContext(r.Context(), in.Event, []byte(in.Body))
	if err != nil && len(results) == 0 {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": results, "error": errString(err)})
}

func (s *Server) handleDeliveries(w http.ResponseWriter, r *http.Request) {
	n := 50
	if q := r.URL.Query().Get("n"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			n = v
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"deliveries": s.hub.RecentDeliveries(n)})
}

func (s *Server) handleStats(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, s.hub.Stats())
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
