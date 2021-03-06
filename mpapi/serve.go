package mpapi

import (
	"net/http"

	"github.com/mjm/courier-js/pkg/tracehttp"

	"github.com/mjm/mpsanity"
	"github.com/mjm/mpsanity/block"
)

type MicropubHandler struct {
	Sanity     *mpsanity.Client
	docBuilder DocumentBuilder
	baseURL    string
	webhookURL string
	authz      Authorizer
	mux        *http.ServeMux
}

func New(sanity *mpsanity.Client, opts ...Option) *MicropubHandler {
	h := &MicropubHandler{
		Sanity: sanity,
		docBuilder: &DefaultDocumentBuilder{
			MarkdownConverter: block.NewMarkdownConverter(),
		},
		mux: http.NewServeMux(),
	}

	for _, o := range opts {
		o.Apply(h)
	}

	h.mux.Handle("/micropub/media", h.AuthMiddleware(http.HandlerFunc(h.handleMedia)))
	h.mux.HandleFunc("/micropub", h.handleMicropub)
	return h
}

func (h *MicropubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "micropub.Request")
	defer span.End()

	r = r.WithContext(ctx)
	tracehttp.WrapHandler(h.mux).ServeHTTP(w, r)
}

type Option interface {
	Apply(h *MicropubHandler)
}

type optionFn func(h *MicropubHandler)

func (f optionFn) Apply(h *MicropubHandler) { f(h) }

func WithDocumentBuilder(b DocumentBuilder) Option {
	return optionFn(func(h *MicropubHandler) {
		h.docBuilder = b
	})
}

func WithBaseURL(u string) Option {
	return optionFn(func(h *MicropubHandler) {
		h.baseURL = u
	})
}

func WithWebhookURL(u string) Option {
	return optionFn(func(h *MicropubHandler) {
		h.webhookURL = u
	})
}

func WithIndieAuth(tokenEndpoint string, me string) Option {
	return optionFn(func(h *MicropubHandler) {
		if tokenEndpoint != "" {
			h.authz = &IndieAuthAuthorizer{
				TokenEndpoint: tokenEndpoint,
				Me:            me,
				HTTPClient:    h.Sanity.HTTPClient,
			}
		}
	})
}
