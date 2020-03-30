package mpsanity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/api/trace"
)

func (c *Client) Txn() *Txn {
	return &Txn{client: c}
}

type Txn struct {
	client    *Client
	mutations []mutation
}

type mutation struct {
	Create            interface{} `json:"create,omitempty"`
	CreateOrReplace   interface{} `json:"createOrReplace,omitempty"`
	CreateIfNotExists interface{} `json:"createIfNotExists,omitempty"`
	Delete            *deletion   `json:"delete,omitempty"`
	// TODO patch
}

type deletion struct {
	ID    string `json:"id,omitempty"`
	Query string `json:"query,omitempty"`
}

func (t *Txn) Create(doc interface{}) *Txn {
	t.mutations = append(t.mutations, mutation{
		Create: doc,
	})
	return t
}

func (t *Txn) CreateOrReplace(doc interface{}) *Txn {
	t.mutations = append(t.mutations, mutation{
		CreateOrReplace: doc,
	})
	return t
}

func (t *Txn) CreateIfNotExists(doc interface{}) *Txn {
	t.mutations = append(t.mutations, mutation{
		CreateIfNotExists: doc,
	})
	return t
}

func (t *Txn) Delete(id string) *Txn {
	t.mutations = append(t.mutations, mutation{
		Delete: &deletion{ID: id},
	})
	return t
}

func (t *Txn) Commit(ctx context.Context) error {
	c := t.client
	ctx, span := tracer.Start(ctx, "txn.Commit",
		trace.WithAttributes(
			projectIDKey(c.ProjectID),
			datasetKey(c.Dataset),
			mutationCountKey(len(t.mutations))))
	defer span.End()

	m := mutationRequest{Mutations: t.mutations}
	body, err := json.Marshal(m)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	r, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/data/mutate/%s", c.Dataset), bytes.NewReader(body))
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	res, err := c.HTTPClient.Do(r)
	if err != nil {
		span.RecordError(ctx, err)
		return err
	}

	if res.StatusCode >= 500 {
		// TODO parse messages out of error response
		err = fmt.Errorf("unexpected server error %d", res.StatusCode)
		span.RecordError(ctx, err)
		return err
	}

	if res.StatusCode >= 400 {
		// TODO parse messages out of error response
		err = fmt.Errorf("unexpected client error %d", res.StatusCode)
		span.RecordError(ctx, err)
		return err
	}

	return nil
}

type mutationRequest struct {
	Mutations []mutation `json:"mutations"`
}
