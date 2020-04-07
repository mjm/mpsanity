package mpapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
)

func (h *MicropubHandler) notifyWebhook(ctx context.Context, title string) {
	if h.webhookURL == "" {
		return
	}

	ctx, span := tracer.Start(ctx, "notifyWebhook",
		trace.WithAttributes(key.String("notify.title", title)))
	defer span.End()

	q := url.Values{
		"trigger_title": []string{title},
	}
	u := h.webhookURL + "?" + q.Encode()

	span.SetAttributes(key.String("notify.url", u))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader("{}"))
	if err != nil {
		span.RecordError(ctx, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := h.Sanity.HTTPClient.Do(req)
	if err != nil {
		span.RecordError(ctx, err)
		return
	}

	if res.StatusCode > 299 {
		span.RecordError(ctx, fmt.Errorf("unexpected status code %d for webhook", res.StatusCode))
		return
	}
}
