package mpsanity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

func (c *Client) UploadImage(ctx context.Context, body io.Reader) (string, error) {
	ctx, span := tracer.Start(ctx, "sanity.Upload",
		trace.WithAttributes(
			projectIDKey(c.ProjectID),
			datasetKey(c.Dataset)))
	defer span.End()

	r, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/assets/images/%s", c.Dataset), body)
	if err != nil {
		span.RecordError(ctx, err)
		return "", err
	}

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		span.RecordError(ctx, err)
		return "", err
	}

	if res.StatusCode >= 500 {
		// TODO parse messages out of error response
		err = fmt.Errorf("unexpected server error %d", res.StatusCode)
		span.RecordError(ctx, err)
		return "", err
	}

	if res.StatusCode >= 400 {
		// TODO parse messages out of error response
		err = fmt.Errorf("unexpected client error %d", res.StatusCode)
		span.RecordError(ctx, err)
		return "", err
	}

	var idResp struct {
		Document struct {
			ID string `json:"_id"`
		} `json:"document"`
	}
	if err := json.NewDecoder(res.Body).Decode(&idResp); err != nil {
		span.RecordError(ctx, err)
		return "", err
	}

	span.SetAttributes(docIDKey(idResp.Document.ID))
	return idResp.Document.ID, nil
}
