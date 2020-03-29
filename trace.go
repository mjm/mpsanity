package mpsanity

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
)

var tracer = global.Tracer("github.com/mjm/mpsanity")

var docIDKey = key.New("sanity.doc_id").String
