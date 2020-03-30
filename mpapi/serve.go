package mpapi

import (
	"net/http"

	"github.com/mjm/courier-js/pkg/tracehttp"

	"github.com/mjm/mpsanity"
)

type MicropubHandler struct {
	Sanity *mpsanity.Client
	mux    *http.ServeMux
}

func New(sanity *mpsanity.Client) *MicropubHandler {
	h := &MicropubHandler{
		Sanity: sanity,
		mux:    http.NewServeMux(),
	}
	h.mux.HandleFunc("/micropub/media", h.handleMedia)
	h.mux.HandleFunc("/micropub", h.handleMicropub)
	return h
}

func (h *MicropubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "micropub.Request")
	defer span.End()

	r = r.WithContext(ctx)
	tracehttp.WrapHandler(h.mux).ServeHTTP(w, r)
}
