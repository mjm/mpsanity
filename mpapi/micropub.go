package mpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

func (h *MicropubHandler) handleMicropub(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleMicropubGet(w, r)
	case http.MethodPost:
		h.handleMicropubPost(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MicropubHandler) handleMicropubGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	q := r.FormValue("q")

	switch q {
	case "config":
		w.Header().Set("Content-Type", "application/json")
		cfg := MicropubConfig{
			MediaEndpoint: fmt.Sprintf("https://%s/micropub/media", r.Host),
		}
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			span.RecordError(ctx, err)
			return
		}
	}
}

func (h *MicropubHandler) handleMicropubPost(w http.ResponseWriter, r *http.Request) {

}
