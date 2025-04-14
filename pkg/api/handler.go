package api

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	limiter  LimiterInterface
	requests *prometheus.CounterVec
}

type LimiterInterface interface {
	Allow(clientID string) bool
}

func (h *Handler) CheckRateLimit(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client")
	if clientID == "" {
		http.Error(w, "client parameter is required", http.StatusBadRequest)
		return
	}

	allowed := h.limiter.Allow(clientID)
	status := "denied"
	if allowed {
		status = "allowed"
	}
	h.requests.WithLabelValues(clientID, status).Inc()

	resp := map[string]bool{"allowed": allowed}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) MetricsHandler() http.Handler {
	return promhttp.Handler()
}
