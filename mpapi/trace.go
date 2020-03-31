package mpapi

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
)

var tracer = global.Tracer("github.com/mjm/mpsanity/mpapi")

var (
	contentLengthKey = key.New("http.content_length").Int

	requestTypeKey    = key.New("micropub.request_type").String
	requestTypeCreate = requestTypeKey("create")
	requestTypeUpdate = requestTypeKey("update")

	typeKey = key.New("micropub.type").String
)
