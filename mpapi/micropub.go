package mpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrMethodNotAllowed = errors.New("method not allowed")

func (h *MicropubHandler) handleMicropub(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	switch r.Method {
	case http.MethodGet:
		h.handleMicropubGet(w, r)
	case http.MethodPost:
		h.handleMicropubPost(w, r)
	default:
		span.RecordError(ctx, ErrMethodNotAllowed)
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
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
			respondWithError(ctx, w, err)
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
		respondWithError(ctx, w, status.Error(codes.InvalidArgument, err.Error()))
		return
	}

	switch mediaType {
	case "application/json":
		h.handleMicropubJSON(w, r)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		h.handleMicropubForm(w, r)
	}
}

func (h *MicropubHandler) createDocument(ctx context.Context, w http.ResponseWriter, input *CreateInput) {
	span := trace.SpanFromContext(ctx)

	doc, err := h.docBuilder.BuildDocument(ctx, input)
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	if err := h.Sanity.Txn().Create(doc).Commit(ctx); err != nil {
		respondWithError(ctx, w, err)
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
