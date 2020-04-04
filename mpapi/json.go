package mpapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/mpsanity"
)

func (h *MicropubHandler) handleMicropubJSON(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleMicropubJSON")
	defer span.End()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	span.SetAttributes(contentLengthKey(len(data)))

	var typeVal struct {
		Action string   `json:"action"`
		Type   []string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeVal); err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		respondWithError(ctx, w, err)
		return
	}

	if typeVal.Action == "update" {
		span.SetAttributes(requestTypeUpdate)
		// TODO handle update
	} else if len(typeVal.Type) > 0 {
		// create
		span.SetAttributes(requestTypeCreate, typeKey(typeVal.Type[0]))

		var input CreateInput
		if err := json.Unmarshal(data, &input); err != nil {
			err = status.Error(codes.InvalidArgument, err.Error())
			respondWithError(ctx, w, err)
			return
		}

		input.Type[0] = strings.TrimPrefix(input.Type[0], "h-")

		if len(input.Props.Photo) > 0 {
			rs, err := h.fetchImageAssets(ctx, input.Props.Photo)
			if err != nil {
				respondWithError(ctx, w, err)
				return
			}

			imgIDs, err := h.uploadImageAssets(ctx, rs)
			if err != nil {
				respondWithError(ctx, w, err)
				return
			}

			input.Props.Photo = imgIDs
		}

		doc, err := h.docBuilder.BuildDocument(ctx, &input)
		if err != nil {
			respondWithError(ctx, w, err)
			return
		}

		h.createDocument(ctx, w, doc)
	}
}

type CreateInput struct {
	Type  []string `json:"type"`
	Props struct {
		Name        []string    `json:"name,omitempty"`
		Content     []string    `json:"content,omitempty"`
		Slug        []string    `json:"mp-slug"`
		Published   []time.Time `json:"published"`
		Photo       []string    `json:"photo"`
		Syndication []string    `json:"syndication"`
	} `json:"properties"`
}

func (in *CreateInput) Name() string {
	if vs := in.Props.Name; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Content() string {
	if vs := in.Props.Content; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Slug() string {
	if vs := in.Props.Slug; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Published() *time.Time {
	if vs := in.Props.Published; len(vs) > 0 {
		return &vs[0]
	}
	return nil
}

func (in *CreateInput) Photos() []mpsanity.Reference {
	var refs []mpsanity.Reference
	for _, p := range in.Props.Photo {
		refs = append(refs, mpsanity.Reference(p))
	}
	return refs
}

func (in *CreateInput) Syndication() []string {
	return in.Props.Syndication
}
