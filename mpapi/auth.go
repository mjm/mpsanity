package mpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorizer interface {
	Authorize(r *http.Request) (string, []string, error)
}

type IndieAuthAuthorizer struct {
	TokenEndpoint string
	Me            string
	HTTPClient    *http.Client
}

var (
	ErrNoToken   = status.Error(codes.Unauthenticated, "no token provided in Authorization header")
	ErrBadHeader = status.Error(codes.Unauthenticated, "invalid format for Authorization header")
	ErrRejected  = status.Error(codes.PermissionDenied, "token endpoint rejected token")
	ErrForbidden = status.Error(codes.PermissionDenied, "you are not allowed to access this API")
)

func (authz *IndieAuthAuthorizer) Authorize(r *http.Request) (string, []string, error) {
	ctx, span := tracer.Start(r.Context(), "IndieAuthAuthorizer.Authorize",
		trace.WithAttributes(
			key.String("auth.token_endpoint", authz.TokenEndpoint),
			key.String("auth.expected_me", authz.Me)))
	defer span.End()

	u, err := url.Parse(authz.Me)
	if err != nil {
		span.RecordError(ctx, err)
		return "", nil, err
	}

	expectedMe := u.Hostname()

	authzHeader := r.Header.Get("Authorization")
	if authzHeader == "" {
		span.RecordError(ctx, ErrNoToken, trace.WithErrorStatus(status.Code(ErrNoToken)))
		return "", nil, ErrNoToken
	}
	if !strings.HasPrefix(authzHeader, "Bearer ") {
		span.RecordError(ctx, ErrBadHeader, trace.WithErrorStatus(status.Code(ErrBadHeader)))
		return "", nil, ErrBadHeader
	}

	token := strings.TrimPrefix(authzHeader, "Bearer ")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, authz.TokenEndpoint, nil)
	if err != nil {
		span.RecordError(ctx, err)
		return "", nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := authz.HTTPClient.Do(req)
	if err != nil {
		span.RecordError(ctx, err)
		return "", nil, err
	}

	if res.StatusCode == http.StatusForbidden {
		span.RecordError(ctx, ErrRejected, trace.WithErrorStatus(status.Code(ErrRejected)))
		return "", nil, ErrRejected
	}

	if res.StatusCode > 299 {
		err := status.Errorf(codes.Internal, "bad response from token endpoint: %d", res.StatusCode)
		span.RecordError(ctx, err, trace.WithErrorStatus(status.Code(err)))
		return "", nil, err
	}

	var response struct {
		Me    string `json:"me"`
		Scope string `json:"scope"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		span.RecordError(ctx, err)
		return "", nil, err
	}

	span.SetAttributes(
		key.String("auth.actual_me", response.Me),
		key.String("auth.scope", response.Scope))

	u, err = url.Parse(response.Me)
	if err != nil {
		span.RecordError(ctx, err)
		return "", nil, err
	}

	actualMe := u.Hostname()

	if expectedMe != actualMe {
		span.RecordError(ctx, ErrForbidden, trace.WithErrorStatus(status.Code(ErrForbidden)))
		return "", nil, ErrForbidden
	}

	return actualMe, strings.Split(response.Scope, " "), nil
}

func (h *MicropubHandler) AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.authz == nil {
			handler.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		me, scope, err := h.authz.Authorize(r)
		if err != nil {
			respondWithError(ctx, w, err)
			return
		}

		span.SetAttributes(
			key.String("user", me),
			key.String("auth.scope", strings.Join(scope, " ")))

		handler.ServeHTTP(w, r)
	})
}
