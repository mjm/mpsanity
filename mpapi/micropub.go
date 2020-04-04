package mpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"

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

func (h *MicropubHandler) createDocument(ctx context.Context, w http.ResponseWriter, doc Document) {
	span := trace.SpanFromContext(ctx)

	if err := h.Sanity.Txn().Create(doc).Commit(ctx); err != nil {
		span.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.webhookURL != "" {
		q := url.Values{
			"trigger_title": []string{fmt.Sprintf("Create %s", doc.URLPath())},
		}
		u := h.webhookURL + "?" + q.Encode()
		res, err := http.Post(u, "application/json", strings.NewReader("{}"))
		if err != nil {
			span.RecordError(ctx, err)
		} else if res.StatusCode > 299 {
			span.RecordError(ctx, fmt.Errorf("unexpected status code %d for webhook", res.StatusCode))
		}
	}

	w.Header().Set("Location", h.baseURL+doc.URLPath())
	w.WriteHeader(http.StatusAccepted)
}
