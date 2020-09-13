package mpapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *MicropubHandler) handleMedia(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleMedia")
	defer span.End()

	r.ParseMultipartForm(32 << 20)
	if r.MultipartForm == nil {
		respondWithError(ctx, w, status.Error(codes.InvalidArgument, "media upload requires multipart form data"))
		return
	}

	if len(r.MultipartForm.File["file"]) != 1 {
		respondWithError(ctx, w, status.Error(codes.InvalidArgument, "media upload should have exactly one 'file' field"))
		return
	}

	fh := r.MultipartForm.File["file"][0]
	f, err := fh.Open()
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	imgIDs, err := h.uploadImageAssets(ctx, []io.ReadCloser{f})
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	imgURL := h.urlForImageAsset(imgIDs[0])
	w.Header().Set("Location", imgURL)
	w.WriteHeader(http.StatusCreated)
}

func (h *MicropubHandler) fetchImageAssets(ctx context.Context, urls []string) ([]io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "fetchImageAssets",
		trace.WithAttributes(key.Int("asset_count", len(urls))))
	defer span.End()

	group, subCtx := errgroup.WithContext(ctx)

	rs := make([]io.ReadCloser, len(urls))
	for i, u := range urls {
		group.Go(func() error {
			req, err := http.NewRequestWithContext(subCtx, http.MethodGet, u, nil)
			if err != nil {
				return err
			}

			res, err := h.Sanity.HTTPClient.Do(req)
			if err != nil {
				return err
			}

			if res.StatusCode > 299 {
				return fmt.Errorf("unexpected status code %d for %s", res.StatusCode, u)
			}

			rs[i] = res.Body
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		span.RecordError(ctx, err)
		return nil, err
	}

	return rs, nil
}

func (h *MicropubHandler) uploadImageAssets(ctx context.Context, rs []io.ReadCloser) ([]string, error) {
	ctx, span := tracer.Start(ctx, "uploadImageAssets",
		trace.WithAttributes(key.Int("asset_count", len(rs))))
	defer span.End()

	group, subCtx := errgroup.WithContext(ctx)

	imgIDs := make([]string, len(rs))
	for i, r := range rs {
		group.Go(func() error {
			defer r.Close()

			imgID, err := h.Sanity.UploadImage(subCtx, r)
			if err != nil {
				return err
			}

			imgIDs[i] = imgID
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		span.RecordError(ctx, err)
		return nil, err
	}

	return imgIDs, nil
}

func (h *MicropubHandler) urlForImageAsset(id string) string {
	urlID := strings.TrimPrefix(id, "image-")
	if lastDashIdx := strings.LastIndex(urlID, "-"); lastDashIdx != -1 {
		urlID = urlID[:lastDashIdx] + "." + urlID[lastDashIdx+1:]
	}

	return fmt.Sprintf("https://cdn.sanity.io/images/%s/%s/%s", h.Sanity.ProjectID, h.Sanity.Dataset, urlID)
}
