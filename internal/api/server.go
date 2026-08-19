package api

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/LYH2263/go-webhookhub"
)

type Options struct {
	WebDir    string
	WebFS     fs.FS
	AllowCORS bool
}

type Server struct {
	hub  *webhookhub.Hub
	mux  *http.ServeMux
	opts Options
}

func New(hub *webhookhub.Hub, opts Options) *Server {
	s := &Server{hub: hub, mux: http.NewServeMux(), opts: opts}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler { return s }

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/health", s.handleHealth)
	s.mux.HandleFunc("GET /api/stats", s.handleStats)
	s.mux.HandleFunc("GET /api/endpoints", s.handleListEndpoints)
	s.mux.HandleFunc("POST /api/endpoints", s.handleCreateEndpoint)
	s.mux.HandleFunc("DELETE /api/endpoints/{id}", s.handleDeleteEndpoint)
	s.mux.HandleFunc("POST /api/endpoints/{id}/enable", s.handleEnable)
	s.mux.HandleFunc("POST /api/endpoints/{id}/disable", s.handleDisable)
	s.mux.HandleFunc("POST /api/dispatch", s.handleDispatch)
	s.mux.HandleFunc("GET /api/deliveries", s.handleDeliveries)
	s.mux.HandleFunc("GET /", s.handleStatic)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "ts": time.Now().UTC()})
}

func lookupWeb(dir string) string {
	if dir == "" {
		dir = "web"
	}
	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		return dir
	}
	exe, err := os.Executable()
	if err == nil {
		cand := filepath.Join(filepath.Dir(exe), "web")
		if st, err := os.Stat(cand); err == nil && st.IsDir() {
			return cand
		}
	}
	return dir
}
