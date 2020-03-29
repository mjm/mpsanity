package mpsanity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/api/trace"
)

func (c *Client) Query(ctx context.Context, query string, out interface{}) error {
	ctx, span := tracer.Start(ctx, "sanity.Query",
		trace.WithAttributes(
			projectIDKey(c.ProjectID),
			datasetKey(c.Dataset),
			queryKey(query)))
	defer span.End()

	q := url.Values{
		"query": []string{query},
	}
	r, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/data/query/%s?%s", c.Dataset, q.Encode()), nil)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	var result queryResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		span.RecordError(ctx, err)
		return err
	}

	if result.Result != nil {
		if err := json.Unmarshal(*result.Result, out); err != nil {
			span.RecordError(ctx, err)
			return err
		}
	}

	return nil
}
