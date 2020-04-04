package mpapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"golang.org/x/sync/errgroup"
)

func (h *MicropubHandler) handleMedia(w http.ResponseWriter, r *http.Request) {

}

func (h *MicropubHandler) fetchImageAssets(ctx context.Context, urls []string) ([]io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "fetchImageAssets",
		trace.WithAttributes(key.Int("asset_count", len(urls))))
	defer span.End()

	group, _ := errgroup.WithContext(ctx)

	rs := make([]io.ReadCloser, len(urls))
	for i, u := range urls {
		group.Go(func() error {
			res, err := http.Get(u)
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
