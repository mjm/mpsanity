package mpapi

import (
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer("github.com/mjm/mpsanity/mpapi")
