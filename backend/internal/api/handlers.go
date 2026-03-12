package api

import (
	"encoding/json"
	"net/http"

	"smartops/backend/internal/devicemanager"
)

type Server struct{}

type storagePoolRequest struct {
	devicemanager.Config
}

type lunRequest struct {
	devicemanager.Config
	LUNID string `json:"lunId"`
}

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/storagepools", s.handleStoragePools)
	mux.HandleFunc("/api/lun", s.handleLUN)
	mux.HandleFunc("/api/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}

func (s *Server) handleStoragePools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req storagePoolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	client, err := devicemanager.NewClient(req.Config)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	pools, err := client.GetStoragePools()
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": pools})
}

func (s *Server) handleLUN(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	var req lunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if req.LUNID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "lunId is required"})
		return
	}
	client, err := devicemanager.NewClient(req.Config)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	lun, err := client.GetLUN(req.LUNID)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": lun})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
