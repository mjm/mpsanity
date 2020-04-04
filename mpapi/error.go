package mpapi

import (
	"context"
	"net/http"

	"github.com/mjm/courier-js/pkg/tracehttp"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/status"
)

func respondWithError(ctx context.Context, w http.ResponseWriter, err error) {
	span := trace.SpanFromContext(ctx)
	code := status.Code(err)
	span.RecordError(ctx, err, trace.WithErrorStatus(code))
	http.Error(w, err.Error(), tracehttp.StatusCode(code))
}
