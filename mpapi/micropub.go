package mpapi

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
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
	span.SetAttributes(key.String("micropub.q", q))

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
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	mediaType := r.Header.Get("Content-Type")
	span.SetAttributes(key.String("http.content_type", mediaType))

	mediaType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.InvalidArgument))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch mediaType {
	case "application/json":
		h.handleMicropubJSON(w, r)
	case "application/x-www-form-urlencoded":
		h.handleMicropubURLEncoded(w, r)
	case "multipart/form-data":
		h.handleMicropubMultipart(w, r)
	}
}

func (h *MicropubHandler) handleMicropubURLEncoded(w http.ResponseWriter, r *http.Request) {

}

func (h *MicropubHandler) handleMicropubMultipart(w http.ResponseWriter, r *http.Request) {

}
