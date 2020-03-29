package mpsanity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

type docResult struct {
	Docs []json.RawMessage `json:"documents"`
}

func (c *Client) Doc(ctx context.Context, id string, out interface{}) error {
	ctx, span := tracer.Start(ctx, "sanity.Doc",
		trace.WithAttributes(docIDKey(id)))
	defer span.End()

	r, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/data/doc/%s/%s", c.Dataset, id), nil)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	var result docResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		span.RecordError(ctx, err)
		return err
	}

	if len(result.Docs) > 0 {
		if err := json.Unmarshal(result.Docs[0], out); err != nil {
			span.RecordError(ctx, err)
			return err
		}
	}

	return nil
}
