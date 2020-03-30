package mpsanity

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
)

var tracer = global.Tracer("github.com/mjm/mpsanity")

var (
	projectIDKey     = key.New("sanity.project_id").String
	datasetKey       = key.New("sanity.dataset").String
	docIDKey         = key.New("sanity.doc_id").String
	queryKey         = key.New("sanity.query").String
	mutationCountKey = key.New("sanity.mutation_count").Int
)
