package mpapi

import (
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *MicropubHandler) handleMicropubForm(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleMicropubForm",
		trace.WithAttributes(requestTypeCreate))
	defer span.End()

	r = r.WithContext(ctx)

	input, err := h.readInputFromForm(r)
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	var rs []io.ReadCloser
	for _, fh := range r.MultipartForm.File["photo"] {
		f, err := fh.Open()
		if err != nil {
			respondWithError(ctx, w, err)
			return
		}

		rs = append(rs, f)
	}
	for _, fh := range r.MultipartForm.File["photo[]"] {
		f, err := fh.Open()
		if err != nil {
			respondWithError(ctx, w, err)
			return
		}

		rs = append(rs, f)
	}

	imgIDs, err := h.uploadImageAssets(ctx, rs)
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	input.Props.Photo = imgIDs

	h.createDocument(ctx, w, input)
}

func (h *MicropubHandler) readInputFromForm(r *http.Request) (*CreateInput, error) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)

	r.ParseMultipartForm(32 << 20)

	input := new(CreateInput)
	input.Type = r.Form["h"]
	span.SetAttributes(typeKey(input.Type[0]))
	input.Props.Name = r.Form["name"]
	input.Props.Content = r.Form["content"]
	input.Props.Slug = r.Form["mp-slug"]

	for _, pub := range r.Form["published"] {
		t, err := time.Parse(time.RFC3339, pub)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		input.Props.Published = append(input.Props.Published, t)
	}

	return input, nil
}
